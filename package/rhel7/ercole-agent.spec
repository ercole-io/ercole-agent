Name:           ercole-agent
Version:        ERCOLE_VERSION
Release:        1%{?dist}
Summary:        Agent for ercole

License:        Proprietary
URL:            https://github.com/ercole-io/%{name}
Source0:        https://github.com/ercole-io/%{name}/archive/%{name}-%{version}.tar.gz
Requires: bc systemd
BuildRequires: systemd

Group:          Tools

Buildroot: /tmp/rpm-ercole-agent

%global debug_package %{nil}

%description
Ercole Agent collects information about the Oracle DB instances
running on the local machine and send information to a central server

%pre
getent passwd ercole >/dev/null || \
    useradd -r -d /home/ercole-agent -m -s /bin/bash \
    -c "Ercole agent user" ercole
getent passwd ercole >/dev/null || \
    useradd -r -d /home/ercole-agent -m -s /bin/bash \
    -c "Ercole agent user" ercole
getent passwd ercole >/dev/null || \
    useradd -r -d /home/ercole-agent -m -s /bin/bash \
    -c "Ercole agent user" ercole
getent passwd ercole >/dev/null || \
    useradd -r -d /home/ercole-agent -m -s /bin/bash \
    -c "Ercole agent user" ercole

if getent group oinstall >/dev/null; then
    usermod -aG oinstall ercole
fi

if getent group dba >/dev/null; then
    usermod -aG dba ercole
fi

if getent group mysql >/dev/null; then
    usermod -aG mysql ercole
fi
exit 0

%prep
%setup -q -n %{name}-%{version}

rm -rf %{buildroot}
make DESTDIR=%{buildroot}/opt/ercole-agent install
install -d %{buildroot}/etc/systemd/system
install -d %{buildroot}/opt/ercole-agent/run
install -d %{buildroot}%{_unitdir} 
install -d %{buildroot}%{_presetdir}
install -m 0644 -C package/rhel7/ercole-agent.service %{buildroot}%{_unitdir}/%{name}.service
install -m 0644 package/rhel7/60-ercole-agent.preset %{buildroot}%{_presetdir}/60-%{name}.preset

%post
/usr/bin/systemctl preset %{name}.service >/dev/null 2>&1 ||:
if [ -e /opt/ercole-agent/.dbs_group ]; then
  echo "File already exists. Do not overwrite."
else
  echo "File does not exist. Creating an empty file."
  touch /opt/ercole-agent/.dbs_group
fi
if [ -e /opt/ercole-agent/.cell_group ]; then
  echo "File already exists. Do not overwrite."
else
  echo "File does not exist. Creating an empty file."
  touch /opt/ercole-agent/.cell_group
fi
if [ -e /opt/ercole-agent/.ibs_group_EMPTY ]; then
  echo "File already exists. Do not overwrite."
else
  echo "File does not exist. Creating an empty file."
  touch /opt/ercole-agent/.ibs_group_EMPTY
fi
/usr/bin/systemctl enable %{name}.service >/dev/null 2>&1 || :

%preun
if [ $1 -eq 0 ]; then
  /usr/bin/systemctl --no-reload disable %{name}.service >/dev/null 2>&1 || :
  /usr/bin/systemctl stop %{name}.service >/dev/null 2>&1 ||:
fi

%postun
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:

%files
%attr(-,ercole,-) /opt/ercole-agent/run
%dir /opt/ercole-agent
%dir /opt/ercole-agent/fetch
%dir /opt/ercole-agent/sql
%config(noreplace) /opt/ercole-agent/config.json
/opt/ercole-agent/ercole-agent
/opt/ercole-agent/ercole-setup

/opt/ercole-agent/fetch/linux/addm.sh
/opt/ercole-agent/fetch/linux/backup.sh
/opt/ercole-agent/fetch/linux/checkpdb.sh
/opt/ercole-agent/fetch/linux/cloud_membership_aws.sh
/opt/ercole-agent/fetch/linux/cluster_membership_status.sh
/opt/ercole-agent/fetch/linux/db.sh
/opt/ercole-agent/fetch/linux/dbmounted.sh
/opt/ercole-agent/fetch/linux/dbstatus.sh
/opt/ercole-agent/fetch/linux/dbversion.sh
/opt/ercole-agent/fetch/linux/filesystem.sh
/opt/ercole-agent/fetch/linux/grant_dba.sh
/opt/ercole-agent/fetch/linux/host.sh
/opt/ercole-agent/fetch/linux/cwversion.sh
/opt/ercole-agent/fetch/linux/license.sh
/opt/ercole-agent/fetch/linux/listpdb.sh
/opt/ercole-agent/fetch/linux/opt.sh
/opt/ercole-agent/fetch/linux/oratab.sh
/opt/ercole-agent/fetch/linux/ovm.sh
/opt/ercole-agent/fetch/linux/patch.sh
/opt/ercole-agent/fetch/linux/psu.sh
/opt/ercole-agent/fetch/linux/schema.sh
/opt/ercole-agent/fetch/linux/schema_pdb.sh
/opt/ercole-agent/fetch/linux/segmentadvisor.sh
/opt/ercole-agent/fetch/linux/segmentadvisor_pdb.sh
/opt/ercole-agent/fetch/linux/services.sh
/opt/ercole-agent/fetch/linux/services_pdb.sh
/opt/ercole-agent/fetch/linux/stats.sh
/opt/ercole-agent/fetch/linux/tablespace.sh
/opt/ercole-agent/fetch/linux/tablespace_pdb.sh
/opt/ercole-agent/fetch/linux/oracle_running_databases.sh
/opt/ercole-agent/fetch/linux/rac.sh
/opt/ercole-agent/fetch/linux/size_pdb.sh
/opt/ercole-agent/fetch/linux/charset_pdb.sh
/opt/ercole-agent/fetch/linux/exec_verbose.sh
/opt/ercole-agent/fetch/linux/partitioning.sh
/opt/ercole-agent/fetch/linux/partitioning_pdb.sh
/opt/ercole-agent/fetch/linux/vmware.ps1
/opt/ercole-agent/fetch/linux/cdb_cpu_iops.sh
/opt/ercole-agent/fetch/linux/pdb_cpu_iops.sh
/opt/ercole-agent/fetch/linux/sar_cpu_only_linux.sh
/opt/ercole-agent/fetch/linux/sar_disks_only_linux.sh
/opt/ercole-agent/fetch/linux/to_postgresql.sh
/opt/ercole-agent/fetch/linux/to_postgresql_pluggable.sh
/opt/ercole-agent/fetch/linux/memory_pga_sga_advisory.sh

/opt/ercole-agent/fetch/linux/exadata/info.sh
/opt/ercole-agent/fetch/linux/exadata/new_info.sh
/opt/ercole-agent/fetch/linux/exadata/storage-status.sh

/opt/ercole-agent/fetch/linux/mysql/mysql_gather.sh

/opt/ercole-agent/fetch/linux/postgresql/psql.sh
/opt/ercole-agent/fetch/linux/postgresql/psql_schema.sh

/opt/ercole-agent/fetch/linux/suggest_oratab.sh
/opt/ercole-agent/fetch/linux/oracle_running_database_home_oratabless.sh
/opt/ercole-agent/fetch/linux/oracle_running_databases_oratabless.sh

/opt/ercole-agent/sql/addm.sql
/opt/ercole-agent/sql/backup_schedule.sql
/opt/ercole-agent/sql/checkpdb.sql
/opt/ercole-agent/sql/db.sql
/opt/ercole-agent/sql/dbmounted.sql
/opt/ercole-agent/sql/dbone.sql
/opt/ercole-agent/sql/dbstatus.sql
/opt/ercole-agent/sql/edition.sql
/opt/ercole-agent/sql/grant_dba.sql
/opt/ercole-agent/sql/license-10.sql
/opt/ercole-agent/sql/license.sql
/opt/ercole-agent/sql/license_pluggable.sql
/opt/ercole-agent/sql/listpdb.sql
/opt/ercole-agent/sql/opt.sql
/opt/ercole-agent/sql/opt_pluggable.sql
/opt/ercole-agent/sql/patch-12.sql
/opt/ercole-agent/sql/patch.sql
/opt/ercole-agent/sql/pluggable.sql
/opt/ercole-agent/sql/psu-1.sql
/opt/ercole-agent/sql/psu-2.sql
/opt/ercole-agent/sql/schema.sql
/opt/ercole-agent/sql/schema_pdb.sql
/opt/ercole-agent/sql/segment_advisor.sql
/opt/ercole-agent/sql/segment_advisor_pdb.sql
/opt/ercole-agent/sql/services.sql
/opt/ercole-agent/sql/services_pdb.sql
/opt/ercole-agent/sql/stats.sql
/opt/ercole-agent/sql/ts.sql
/opt/ercole-agent/sql/ts_pdb.sql
/opt/ercole-agent/sql/size_pdb.sql
/opt/ercole-agent/sql/charset_pdb.sql
/opt/ercole-agent/sql/partitioning.sql
/opt/ercole-agent/sql/partitioning_pdb.sql
/opt/ercole-agent/sql/pdb_cpu_iops.sql
/opt/ercole-agent/sql/cdb_cpu_iops.sql
/opt/ercole-agent/sql/to_postgresql.sql
/opt/ercole-agent/sql/memory_pga_sga_advisory.sql

/opt/ercole-agent/sql/mssqlserver/mssqlserver.backup_schedule.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.db.10.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.db.14.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.dbmounted.10.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.dbmounted.14.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.dbstatus.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.edition.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.instanceVersion.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.licensingInfo.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.listDatabases.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.psu-1.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.schema.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.segment_advisor.sql
/opt/ercole-agent/sql/mssqlserver/mssqlserver.ts.sql

/opt/ercole-agent/sql/mysql/databases.sql
/opt/ercole-agent/sql/mysql/high_availability.sql
/opt/ercole-agent/sql/mysql/instance.sql
/opt/ercole-agent/sql/mysql/old_instance.sql
/opt/ercole-agent/sql/mysql/version.sql
/opt/ercole-agent/sql/mysql/segment_advisors.sql
/opt/ercole-agent/sql/mysql/slave_hosts.sql
/opt/ercole-agent/sql/mysql/slave_status.sql
/opt/ercole-agent/sql/mysql/table_schemas.sql

/opt/ercole-agent/sql/postgresql/d_info.sql
/opt/ercole-agent/sql/postgresql/d_info_v10.sql
/opt/ercole-agent/sql/postgresql/erc_GetDB.sql
/opt/ercole-agent/sql/postgresql/erc_GetSchema.sql
/opt/ercole-agent/sql/postgresql/i_info.sql
/opt/ercole-agent/sql/postgresql/i_info_v10.sql
/opt/ercole-agent/sql/postgresql/i_settings.sql
/opt/ercole-agent/sql/postgresql/n_info.sql
/opt/ercole-agent/sql/postgresql/n_info_v10.sql

%config(noreplace) %{_unitdir}/ercole-agent.service
%{_presetdir}/60-ercole-agent.preset

%changelog
* Mon May  7 2018 Simone Rota <srota2@sorint.it>
- 
