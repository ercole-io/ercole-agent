#!/bin/bash
# Author: Attilio Seghezzi
# Date: 31st March 2023
# 20230403: added X9M support and automatic platform detection. Added also VM check function
# 20230418: added check on configuration files' content to ensure they're not empty
# 20230505: enhanced cell entities and added RACK_ID to all physical components
# 20230601: corrected function HostGetDetails to retrieve the correct amount of memory from dom0 type hosts
# 20230727: fixed function CellGetDetails to display the correct storage server name in despite of the network alias used by the utility
# 20231003: added configuration to allow the script to be executed as non-root user. Added summary function fullRun. Added field 'HOST_ID' for both cells and dbnodes

### Variables
APP_DIR=/tmp
ERCOLE_DIR=/opt/ercole-agent
DBS_LST=${ERCOLE_DIR}/.dbs_group
CELL_LST=${ERCOLE_DIR}/.cell_group
IBS_LST=${ERCOLE_DIR}/.ibs_group

### Utility functions
function checkRoot {
    export CURRENT_USR=$(whoami)
        if [ "$CURRENT_USR" != "root" -a "$CURRENT_USR" != "$NONROOT" ]; then
            echo "--> ERROR: utility cannot be executed as a user different then root/$NONROOT (current user: ${CURRENT_USR}). Please switch to a supported one and try again"
            echo "Exiting..."
            exit 1
        fi
}

function checkFiles {
    DIR=$ERCOLE_DIR
    if [[ ! -d "$DIR" ]]; then
            echo "--> ERROR: directory $ERCOLE_DIR does not exist"
            echo "           Please create it and make sure the required configuration files are in it"
            echo "Exiting..."
            exit 1
        fi
}

function checkPWDless {
    if [[ ! -f $1 ]]; then
        echo " --> ERROR: configuration file $1 does not exists. Please check why it is missing and retry"
        echo "Exiting..."
        exit 1
    fi
    CHECK_WC=$(cat $1 |wc -l)
    if [[ "$CHECK_WC" == "0" ]]; then
        echo " --> ERROR: configuration file $1 exists but does not contain any data. Please check its content"
        echo "Exiting..."
        exit 1
    fi 
    LST=$(cat $1)
    for HOST in $LST
    do
        ssh -o BatchMode=yes $HOST -l $2 'exit'
        RC=$?
        if [[ "$RC" != 0 ]]; then
            echo " --> ERROR: passwordless access not set for user ${CURRENT_USR} towards host ${HOST}"
            echo "            Please retry after setting it up using the following command:"
            echo "               dcli -g $1 -l $2 -k"
            echo "Exiting..."
            exit 1
        fi
    done
}

function getRackID {
    export RACK_ID=$(sudo ipmitool sunoem cli 'show /SP system_identifier'|grep '='|awk '{ print $7 }')
}

function checkVM {
    MANUFACTURER=$(sudo dmidecode -s system-manufacturer)
    if [[ "$MANUFACTURER" == "Oracle Corporation" ]]; then
        vm_maker --list-domains 2> /dev/null 1> /dev/null
        RC_VMMAKER=$?
        if [[ "$RC_VMMAKER" != "0" ]]; then
            which xm 2> /dev/null 1> /dev/null
            RC_XM=$?
            if [[ "$RC_XM" == "0" ]]; then
                export PRERUN=dom0
                getRackID
            else
                export PRERUN=bm
                getRackID
            fi
        else 
            export PRERUN=kvm
            getRackID
        fi
    else
        export PRERUN=vm
    fi
}

function preRunFunction {
    checkRoot
    checkFiles
    checkVM
    case $PRERUN in 
        kvm)
            checkPWDless $DBS_LST $NONROOT
            checkPWDless $CELL_LST root
        ;;
        dom0)
            checkPWDless $DBS_LST $NONROOT
            checkPWDless $CELL_LST root
        ;;
        bm)
            checkPWDless $DBS_LST $NONROOT
            checkPWDless $CELL_LST root
        ;;
        vm)
            checkPWDless $DBS_LST $NONROOT
        ;;
        *)
            echo " --> ERROR: the parameter passed is wrong ($PRERUN). Please check if everything is configured as it should"
            echo " --> Exiting..."
            exit 1;
        ;;
        esac
        if [[ -f $IBS_LST ]]; then
            export IB=true
            checkPWDless $IBS_LST root
        fi
}

### Physical nodes functions
function HostGetDetails {
    echo "HOST_TYPE|||RACK_ID|||HOSTNAME|||HOST_ID|||CPU_ENABLED|||CPU_TOT|||MEMORY_GB|||IMAGEVERSION|||KERNEL|||MODEL|||FAN_USED|||FAN_TOTAL|||PSU_USED|||PSU_TOTAL|||MS_STATUS|||RS_STATUS"
    while read NHOST
    do
        INFO=$(dcli -c $NHOST -l $NONROOT "sudo dbmcli -e list dbserver attributes name,cpuCount,fanCount,kernelVersion,powerCount,releaseVersion,msStatus,rsStatus,id"| sed "s/${NHOST}: //g")
        HOST=$(echo $INFO|awk '{print $1}')
        CPU_ENABLED=$(echo $INFO|awk '{print $2}'|awk -F'[/]' '{print $1}')
        CPU_TOT=$(echo $INFO|awk '{print $2}'|awk -F'[/]' '{print $2}')
        FAN_USED=$(echo $INFO|awk '{print $3}'|awk -F'[/]' '{print $1}')
        FAN_TOTAL=$(echo $INFO|awk '{print $3}'|awk -F'[/]' '{print $2}')
        KERNEL=$(echo $INFO|awk '{print $4}')
        PSU_USED=$(echo $INFO|awk '{print $5}'|awk -F'[/]' '{print $1}')
        PSU_TOTAL=$(echo $INFO|awk '{print $5}'|awk -F'[/]' '{print $2}')
        IMAGEVERSION=$(echo $INFO|awk '{print $6}')
        MS=$(echo $INFO|awk '{print $7}')
        RS=$(echo $INFO|awk '{print $8}')
        HOSTID=$(echo $INFO|awk '{print $9}')
        MODEL_APP=$(dcli -c $NHOST -l $NONROOT "sudo ipmitool sunoem cli 'show /SYSTEM system_identifier'|grep system_identifier|grep -iv show"|sed "s/${NHOST}: //g"|awk '{print $3 " " $4 " " $5 " " $6}')
        MODEL_TMP=$(dcli -c $NHOST -l $NONROOT "sudo ipmitool sunoem cli 'show /SYSTEM component_model'|grep component_model|grep -iv show"|sed "s/${NHOST}: //g"|awk '{print $5}')
        if [[ "$MODEL_APP" == "Exadata Database Machine X8M-2" && "$MODEL_TMP" == "X9-2" ]]; then
            MODEL="Exadata Database Machine X9M-2"
        else
            MODEL=$MODEL_APP
        fi
        if [[ "$HOST_TYPE" == "DOM0" ]]; then
            MEMORY_MB=$(dcli -c $NHOST -l $NONROOT "sudo xm info|grep total_memory"|awk '{print $4}')
            MEMORY_GB=$(expr $MEMORY_MB / 1024)
        else
            MEMORY_KB=$(dcli -c $NHOST -l $NONROOT "sudo cat /proc/meminfo|grep MemTotal"|awk '{print $3}')
            MEMORY_GB=$(expr $MEMORY_KB / 1048576)
        fi               
        echo "$RACK_ID|||$HOST_TYPE|||$HOST|||$HOSTID|||$CPU_ENABLED|||$CPU_TOT|||$MEMORY_GB|||$IMAGEVERSION|||$KERNEL|||$MODEL|||$FAN_USED|||$FAN_TOTAL|||$PSU_USED|||$PSU_TOTAL|||$MS|||$RS"
    done < $DBS_LST
}

function BMGetDetails {
    export HOST_TYPE=BARE_METAL
    HostGetDetails
}

function KVMHostGetDetails {
    export HOST_TYPE=KVM_HOST
    HostGetDetails
}

function dom0GetDetails {
    export HOST_TYPE=DOM0
    HostGetDetails
}


### Storage cells functions
function CellGetDetails {
    HOST_TYPE=STORAGE_CELL
    echo "HOST_TYPE|||RACK_ID|||HOST|||HOST_ID|||CPU_ENABLED|||CPU_TOT|||MEMORY_GB|||IMAGEVERSION|||KERNEL|||MODEL|||FAN_USED|||FAN_TOTAL|||PSU_USED|||PSU_TOTAL|||CELLSRV_STATUS|||MS_STATUS|||RS_STATUS"
    while read CELL 
    do
        INFO=$(dcli -c $CELL -l root "cellcli -e list cell attributes name,cpuCount,fanCount,kernelVersion,powerCount,releaseVersion,memoryGB,cellsrvStatus,msStatus"| sed "s/${CELL}: //g")
        DBS=$(dcli -c $CELL -l root "cellcli -e list database attributes name"|sed "s/${CELL}: //g"|sort)
        CDISKS=$(dcli -c $CELL -l root "cellcli -e list celldisk attributes name where diskType=HardDisk"|sed "s/${CELL}: //g"|sort)
        INFO_APP=$(dcli -c $CELL -l root "cellcli -e list cell attributes id,rsStatus"| sed "s/${CELL}: //g")
        MODEL=$(dcli -c $CELL -l root "cellcli -e list cell attributes makeModel"| sed "s/${CELL}: //g")
        HOST=$(echo $INFO|awk '{print $1}')
        CPU_ENABLED=$(echo $INFO|awk '{print $2}'|awk -F'[/]' '{print $1}')
        CPU_TOT=$(echo $INFO|awk '{print $2}'|awk -F'[/]' '{print $2}')
        FAN_USED=$(echo $INFO|awk '{print $3}'|awk -F'[/]' '{print $1}')
        FAN_TOTAL=$(echo $INFO|awk '{print $3}'|awk -F'[/]' '{print $2}')
        KERNEL=$(echo $INFO|awk '{print $4}')
        PSU_USED=$(echo $INFO|awk '{print $5}'|awk -F'[/]' '{print $1}')
        PSU_TOTAL=$(echo $INFO|awk '{print $5}'|awk -F'[/]' '{print $2}')
        IMAGEVERSION=$(echo $INFO|awk '{print $6}')
        MEMORY_GB=$(echo $INFO|awk '{print $7}')
        CELLSRV=$(echo $INFO|awk '{print $8}')
        MS=$(echo $INFO|awk '{print $9}')
        HOSTID=$(echo $INFO_APP|awk '{print $1}')
        RS=$(echo $INFO_APP|awk '{print $2}')
        echo "$HOST_TYPE|||$RACK_ID|||$HOST|||$HOSTID|||$CPU_ENABLED|||$CPU_TOT|||$MEMORY_GB|||$IMAGEVERSION|||$KERNEL|||$MODEL|||$FAN_USED|||$FAN_TOTAL|||$PSU_USED|||$PSU_TOTAL|||$CELLSRV|||$MS|||$RS"
        echo " "; echo "TYPE|||CELLDISK|||CELL|||SIZE|||FREESPACE|||STATUS|||ERROR_COUNT"
        for CDISK in $CDISKS
        do
            TYPE=CELLDISK
            CD_INFO=$(dcli -c $CELL -l root "cellcli -e list celldisk attributes name,errorCount,freeSpace,size,status where name=$CDISK"|sed "s/${CELL}: //g"|sort)
            GDISKS=$(dcli -c $CELL -l root "cellcli -e list griddisk attributes name where cellDisk=$CDISK"|sed "s/${CELL}: //g"|sort)
            CD_NAME=$(echo $CD_INFO|awk '{print $1}')
            CD_ERRCOUNT=$(echo $CD_INFO|awk '{print $2}')
            CD_FREESPACE=$(echo $CD_INFO|awk '{print $3}')
            CD_SIZE=$(echo $CD_INFO|awk '{print $4}')
            CD_STATUS=$(echo $CD_INFO|awk '{print $5}')
            echo "$TYPE|||$CD_NAME|||$HOST|||$CD_SIZE|||$CD_FREESPACE|||$CD_STATUS|||$CD_ERRCOUNT"
            echo "TYPE|||GRIDDISK|||CELLDISK|||SIZE|||STATUS|||ERROR_COUNT|||CACHING_POLICY|||ASMDISK_NAME|||ASM_DIKSKGROUP|||ASKDISK_SIZE|||ASMDISK_STATUS"
            for GDISK in $GDISKS
            do
                TYPE=GRIDDISK
                GD_INFO=$(dcli -c $CELL -l root "cellcli -e list griddisk attributes name,asmDiskGroupName,asmDiskName,asmDiskSize,asmModeStatus,cachingPolicy,errorCount,size,status where name=$GDISK"|sed "s/${CELL}: //g"|sort)
                GD_NAME=$(echo $GD_INFO|awk '{print $1}')
                GD_ASMDGNAME=$(echo $GD_INFO|awk '{print $2}')
                GD_ASMDISKNAME=$(echo $GD_INFO|awk '{print $3}')
                GD_ASMDISKSIZE=$(echo $GD_INFO|awk '{print $4}')
                GD_ASMSTATUS=$(echo $GD_INFO|awk '{print $5}')
                GD_FCPOL=$(echo $GD_INFO|awk '{print $6}')
                GD_ERRCOUNT=$(echo $GD_INFO|awk '{print $7}')
                GD_SIZE=$(echo $GD_INFO|awk '{print $8}')
                GD_STATUS=$(echo $GD_INFO|awk '{print $9}')
                echo "$TYPE|||$GD_NAME|||$CDISK|||$GD_SIZE|||$GD_STATUS|||$GD_ERRCOUNT|||$GD_FCPOL|||$GD_ASMDISKNAME|||$GD_ASMDGNAME|||$GD_ASMDISKSIZE|||$GD_ASMSTATUS"
            done; echo " "
        done
        echo " ";echo "TYPE|||DB_NAME|||CELL|||DBID|||FLASHCACHE_LIMIT|||IORM_SHARE|||LAST_IO_REQ"
        for DB in $DBS 
        do
            TYPE=DATABASE
            DB_INFO=$(dcli -c $CELL -l root "cellcli -e list database attributes name,databaseID,flashCacheLimit,iormShare,lastRequestTime where name=$DB"|sed "s/${CELL}: //g"|sort)
            DB_NAME=$(echo $DB_INFO|awk '{print $1}')
            DB_ID=$(echo $DB_INFO|awk '{print $2}')
            DB_FCLIMIT=$(echo $DB_INFO|awk '{print $3}')
            DB_IORMSHARE=$(echo $DB_INFO|awk '{print $4}')
            DB_LASTREQ=$(echo $DB_INFO|awk '{print $5}')
            echo "$TYPE|||$DB_NAME|||$HOST|||$DB_ID|||$DB_FCLIMIT|||$DB_IORMSHARE|||$DB_LASTREQ"
        done; echo " "
    done < $CELL_LST
}


### VMs functions
function vmGetDetails {
    HOST_TYPE=VM
    echo "HOST_TYPE|||HOSTNAME|||IMAGEVERSION|||KERNEL|||MS_STATUS|||RS_STATUS"
    while read NHOST
    do
        INFO=$(dcli -c $NHOST -l ercole "sudo dbmcli -e list dbserver attributes name,kernelVersion,releaseVersion,msStatus,rsStatus"| sed "s/${NHOST}: //g")
        HOST=$(echo $INFO|awk '{print $1}')
        KERNEL=$(echo $INFO|awk '{print $2}')
        IMAGEVERSION=$(echo $INFO|awk '{print $3}')
        MS=$(echo $INFO|awk '{print $4}')
        RS=$(echo $INFO|awk '{print $5}')
        echo "$HOST_TYPE|||$HOST|||$IMAGEVERSION|||$KERNEL|||$MS|||$RS"
    done < $DBS_LST
}

function KVMGetVMResources {
    HOST_TYPE=VM_KVM
    if [[ "$1" == "-i" ]]; then
        STATUS=Running
    else
        STATUS=NotRunning
    fi
    export RUN_FLAG=$1
    echo "VM_TYPE|||PHYSICAL_HOST|||STATUS|||VM_NAME|||CPU_CURRENT|||CPU_RESTART|||RAM_CURRENT|||RAM_RESTART"
    VM_CPU_LST=${APP_DIR}/.ErcoleAgent_Exa_vCPURunningVMs.lst 
    VM_RAM_LST=${APP_DIR}/.ErcoleAgent_Exa_RAMRunningVMs.lst 
    while read NHOST
    do
        VM_APP=$(dcli -c $NHOST -l $NONROOT sudo vm_maker --list --domain|grep $RUN_FLAG running|awk -F'[(]' '{print $1}')
        VM_LIST=$(echo $VM_APP| sed "s/${NHOST}: //g")
        dcli -g $DBS_LST -l $NONROOT "sudo vm_maker --list --vcpu|grep 'Current:' > $VM_CPU_LST"
        dcli -g $DBS_LST -l $NONROOT "sudo vm_maker --list --memory > $VM_RAM_LST"

        for VM in $VM_LIST
        do
            CPU_CURRENT_APP=$(dcli -c $NHOST -l $NONROOT sudo grep $VM $VM_CPU_LST|awk '{print $4}')
            CPU_CURRENT=$(echo $CPU_CURRENT_APP|sed "s/${NHOST}: //g")
            CPU_RESTART_APP=$(dcli -c $NHOST -l $NONROOT sudo grep $VM $VM_CPU_LST|awk '{print $6}')
            CPU_RESTART=$(echo $CPU_RESTART_APP|sed "s/${NHOST}: //g")
            RAM_CURRENT_APP=$(dcli -c $NHOST -l $NONROOT sudo grep $VM $VM_RAM_LST|awk '{print $3}')
            RAM_CURRENT=$(echo $RAM_CURRENT_APP|sed "s/${NHOST}: //g")
            RAM_RESTART_APP=$(dcli -c $NHOST -l $NONROOT sudo grep $VM $VM_RAM_LST|awk '{print $4}')
            RAM_RESTART=$(echo $RAM_RESTART_APP|sed "s/${NHOST}: //g")
            echo "$HOST_TYPE|||$NHOST|||$STATUS|||$VM|||$CPU_CURRENT|||$CPU_RESTART|||$RAM_CURRENT|||$RAM_RESTART"
        done
    done < $DBS_LST
    dcli -g $DBS_LST -l $NONROOT "sudo rm -f $VM_CPU_LST $VM_RAM_LST"
}

function OVMGetVMResources {
    HOST_TYPE=VM_XEN
    echo "VM_TYPE|||PHYSICAL_HOST|||VM_NAME|||CPU_ONLINE|||CPU_MAX_USABLE|||RAM_ONLINE|||RAM_MAX_USABLE"
    VM_DETAILS=${APP_DIR}/.ErcoleAgent_Exa_VMdetail.lst 
    while read NHOST
    do
        VM_LIST=$(dcli -c $NHOST -l $NONROOT sudo xm list|grep -v Domain-0|grep -v VCPUs|awk '{print $2}')
        for VM in $VM_LIST
        do
            dcli -c $NHOST -l $NONROOT "sudo xm list $VM -l > $VM_DETAILS"
            CPU_ONLINE=$(dcli -c $NHOST -l $NONROOT sudo grep -i online_vcpus $VM_DETAILS|awk '{print $3}'|sed "s/)//g")
            CPU_MAX_USABLE=$(dcli -c $NHOST -l $NONROOT sudo grep -i vcpus $VM_DETAILS|grep -iv online|awk '{print $3}'|sed "s/)//g")
            RAM_ONLINE=$(dcli -c $NHOST -l $NONROOT sudo grep -i memory $VM_DETAILS|grep -iv memory_|grep -iv shadow|awk '{print $3}'|sed "s/)//g")
            RAM_MAX_USABLE=$(dcli -c $NHOST -l $NONROOT sudo grep -i memory_actual $VM_DETAILS|awk '{print $3}'|sed "s/)//g")
            echo "$HOST_TYPE|||$NHOST|||$VM|||$CPU_ONLINE|||$CPU_MAX_USABLE|||$RAM_ONLINE|||$RAM_MAX_USABLE"
        done
    done < $DBS_LST
    dcli -g $DBS_LST -l $NONROOT "sudo rm -f $VM_DETAILS"
}

function KVMGetRunningVMResources {
    KVMGetVMResources -i
}

function KVMGetStoppedVMResources {
    KVMGetVMResources -iv
}


### InfiniBand switches functions
function IBSGetDetails {
    HOST_TYPE=IB_SWITCH
    echo "HOST_TYPE|||RACK_ID|||SWITCH_NAME|||MODEL|||SW_VERSION"
    while read IBS 
    do
        MODEL=$(ibswitches|grep $IBS|awk '{print $6 " " $7 " " $8 " " $9}'|awk -F'["]' '{print $2}')
        VERSION=$(dcli -c $IBS -l root "version|grep SUN"|awk -F'[:]' '{print $3}')
        echo "$HOST_TYPE|||$RACK_ID|||$IBS|||$MODEL|||$VERSION"
    done < $IBS_LST
}


### Summary functions
function KVMFullCheck {
    echo " ";KVMHostGetDetails
    echo " ";KVMGetRunningVMResources
    echo " ";KVMGetStoppedVMResources
}

function dom0FullCheck {
    echo " ";dom0GetDetails
    echo " ";OVMGetVMResources
}

function fullRun {
    preRunFunction
    case $PRERUN in
        kvm)
            KVMFullCheck
            CellGetDetails
        ;;
        dom0)
            dom0FullCheck
            CellGetDetails
            if [[ "$IB" == "true" ]]; then
                IBSGetDetails
            fi
        ;;
        bm)
            BMGetDetails
            CellGetDetails
            if [[ "$IB" == "true" ]]; then
                IBSGetDetails
            fi
        ;;
        vm)
            vmGetDetails
    esac
}


### Main
case $1 in 
    "")
        export NONROOT=root
        fullRun
    ;;
    -cell)
        export NONROOT=root
        getRackID
        checkRoot
        checkFiles
        checkPWDless $CELL_LST
        CellGetDetails
    ;;
    -ibs)
        export NONROOT=root
        getRackID
        checkRoot
        checkFiles
        checkPWDless $IBS_LST
        IBSGetDetails
    ;;
    *)
        export NONROOT=$1
        CHECKUSR=$(grep $NONROOT /etc/passwd|wc -l)
        if [[ "$CHECKUSR" == "1" ]]; then
            CHECKSUDO=$(sudo cat /etc/sudoers|wc -l)
            if [[ "$CHECKSUDO" == "0" ]]; then
                echo " --> ERROR: the specified username exists, but does not have sudo permissions"
                echo "            Please make sure $NONROOT has the correct sudo permissions on all the hosts specified in $DBS_LST"
                echo " --> Exiting..."
                exit 1;
            else
                fullRun
            fi
        else
            echo " --> ERROR: the username/parameter given as input to the shell is wrong ($1)"
            echo "            USERS: please make sure it exists and has sudo permissions"
            echo "            FLAGS: the only supported parameters are -cell/-ibs or blank"
            echo " --> Exiting..."
            exit 1;
        fi
    ;;
esac