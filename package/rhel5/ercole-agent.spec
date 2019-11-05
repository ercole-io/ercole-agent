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
/opt/ercole-agent/fetch/db
/opt/ercole-agent/fetch/dbstatus
/opt/ercole-agent/fetch/feature
/opt/ercole-agent/fetch/opt
/opt/ercole-agent/fetch/filesystem
/opt/ercole-agent/fetch/host
/opt/ercole-agent/fetch/license
/opt/ercole-agent/fetch/oratab
/opt/ercole-agent/fetch/patch
/opt/ercole-agent/fetch/schema
/opt/ercole-agent/fetch/stats
/opt/ercole-agent/fetch/tablespace
/opt/ercole-agent/ercole-agent
/opt/ercole-agent/ercole-setup
/opt/ercole-agent/sql/db.sql
/opt/ercole-agent/sql/feature.sql
/opt/ercole-agent/sql/opt.sql
/opt/ercole-agent/sql/license.sql
/opt/ercole-agent/sql/patch.sql
/opt/ercole-agent/sql/schema.sql
/opt/ercole-agent/sql/stats.sql
/opt/ercole-agent/sql/ts.sql
/etc/init.d/ercole-agent
/etc/logrotate.d/ercole-agent
/opt/ercole-agent/fetch/dbmounted
/opt/ercole-agent/fetch/dbversion
/opt/ercole-agent/sql/dbmounted.sql
/opt/ercole-agent/sql/feature-10.sql
/opt/ercole-agent/sql/license-10.sql
/opt/ercole-agent/sql/patch-12.sql
/opt/ercole-agent/sql/dbstatus.sql
/opt/ercole-agent/sql/edition.sql
/opt/ercole-agent/fetch/addm
/opt/ercole-agent/fetch/psu
/opt/ercole-agent/fetch/segmentadvisor
/opt/ercole-agent/fetch/backup
/opt/ercole-agent/sql/psu-1.sql
/opt/ercole-agent/sql/psu-2.sql
/opt/ercole-agent/sql/addm.sql
/opt/ercole-agent/sql/segment_advisor.sql
/opt/ercole-agent/sql/backup_schedule.sql
/opt/ercole-agent/sql/dbone.sql


%changelog
* Mon May  7 2018 Simone Rota <srota2@sorint.it>
- 
