Name:           ercole-agent
Version:        ERCOLE_VERSION
Release:        1%{?dist}
Summary:        Agent for ercole

License:        Proprietary
URL:            https://wecode.sorint.it/dev.arch/%{name}
Source0:        https://wecode.sorint.it/dev.arch/%{name}/archive/%{name}-%{version}.tar.gz

Group:          Tools
Requires: bc

%define buildroot /tmp/rpm-ercole-agent-buildroot
%define _rpmdir /root

BuildRoot:  %{buildroot}

%global debug_package %{nil}

%description
Ercole Agent collects information about the Oracle DB instances
running on the local machine and send information to a central server

%pre
getent passwd ercole >/dev/null || \
    useradd -r -g oinstall -G oinstall,dba -d /home/ercole-agent -m -s /bin/bash \
    -c "Ercole agent user" ercole
getent passwd ercole >/dev/null || \
    useradd -r -g dba -d /home/ercole-agent -m -s /bin/bash \
    -c "Ercole agent user" ercole
getent passwd ercole >/dev/null || \
    useradd -r -g oinstall -d /home/ercole-agent -m -s /bin/bash \
    -c "Ercole agent user" ercole
exit 0

%prep
%setup -q -n %{name}-%{version}

%build
rm -rf $RPM_BUILD_ROOT
make

%install
make DESTDIR=$RPM_BUILD_ROOT/opt/ercole-agent install
install -d $RPM_BUILD_ROOT/etc/init.d
install -d $RPM_BUILD_ROOT/etc/logrotate.d
install -m 755 package/rhel5/ercole-agent $RPM_BUILD_ROOT/etc/init.d/ercole-agent
install -m 644 package/rhel5/logrotate $RPM_BUILD_ROOT/etc/logrotate.d/ercole-agent

%post

%files
%dir /opt/ercole-agent
%dir /opt/ercole-agent/fetch
%dir /opt/ercole-agent/sql
%config(noreplace) /opt/ercole-agent/config.json
/etc/init.d/ercole-agent
/etc/logrotate.d/ercole-agent
/opt/ercole-agent/ercole-agent
/opt/ercole-agent/ercole-setup
/opt/ercole-agent/fetch/linux/addm.sh
/opt/ercole-agent/fetch/linux/backup.sh
/opt/ercole-agent/fetch/linux/checkpdb.sh
/opt/ercole-agent/fetch/linux/db.sh
/opt/ercole-agent/fetch/linux/dbmounted.sh
/opt/ercole-agent/fetch/linux/dbstatus.sh
/opt/ercole-agent/fetch/linux/dbversion.sh
/opt/ercole-agent/fetch/linux/filesystem.sh
/opt/ercole-agent/fetch/linux/host.sh
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
/opt/ercole-agent/fetch/linux/stats.sh
/opt/ercole-agent/fetch/linux/tablespace.sh
/opt/ercole-agent/fetch/linux/tablespace_pdb.sh
/opt/ercole-agent/fetch/linux/vmware.ps1
/opt/ercole-agent/fetch/mssqlserver/ercoleAgentMsSQLServer-Gather.ps1
/opt/ercole-agent/sql/addm.sql
/opt/ercole-agent/sql/backup_schedule.sql
/opt/ercole-agent/sql/checkpdb.sql
/opt/ercole-agent/sql/db.sql
/opt/ercole-agent/sql/dbmounted.sql
/opt/ercole-agent/sql/dbone.sql
/opt/ercole-agent/sql/dbstatus.sql
/opt/ercole-agent/sql/edition.sql
/opt/ercole-agent/sql/license-10.sql
/opt/ercole-agent/sql/license.sql
/opt/ercole-agent/sql/listpdb.sql
/opt/ercole-agent/sql/opt.sql
/opt/ercole-agent/sql/patch-12.sql
/opt/ercole-agent/sql/patch.sql
/opt/ercole-agent/sql/psu-1.sql
/opt/ercole-agent/sql/psu-2.sql
/opt/ercole-agent/sql/schema.sql
/opt/ercole-agent/sql/schema_pdb.sql
/opt/ercole-agent/sql/segment_advisor.sql
/opt/ercole-agent/sql/stats.sql
/opt/ercole-agent/sql/ts.sql
/opt/ercole-agent/sql/ts_pdb.sql

%changelog
* Mon May  7 2018 Simone Rota <srota2@sorint.it>
- 
