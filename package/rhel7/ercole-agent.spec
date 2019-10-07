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

rm -rf %{buildroot}
make DESTDIR=%{buildroot}/opt/ercole-agent install
install -d %{buildroot}/etc/systemd/system
install -d %{buildroot}/opt/ercole-agent/run
install -d %{buildroot}%{_unitdir} 
install -d %{buildroot}%{_presetdir}
install -m 0644 package/rhel7/ercole-agent.service %{buildroot}%{_unitdir}/%{name}.service
install -m 0644 package/rhel7/60-ercole-agent.preset %{buildroot}%{_presetdir}/60-%{name}.preset

%post
/usr/bin/systemctl preset %{name}.service >/dev/null 2>&1 ||:

%preun
/usr/bin/systemctl --no-reload disable %{name}.service >/dev/null 2>&1 || :
/usr/bin/systemctl stop %{name}.service >/dev/null 2>&1 ||:

%postun
/usr/bin/systemctl daemon-reload >/dev/null 2>&1 ||:

%files
%attr(-,ercole,-) /opt/ercole-agent/run
%dir /opt/ercole-agent
%dir /opt/ercole-agent/fetch
%dir /opt/ercole-agent/sql
%config(noreplace) /opt/ercole-agent/config.json
/opt/ercole-agent/fetch/db
/opt/ercole-agent/fetch/dbstatus
/opt/ercole-agent/fetch/feature
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
/opt/ercole-agent/sql/license.sql
/opt/ercole-agent/sql/patch.sql
/opt/ercole-agent/sql/schema.sql
/opt/ercole-agent/sql/stats.sql
/opt/ercole-agent/sql/ts.sql
/opt/ercole-agent/fetch/dbmounted
/opt/ercole-agent/fetch/dbversion
/opt/ercole-agent/sql/dbmounted.sql
/opt/ercole-agent/sql/feature-10.sql
/opt/ercole-agent/sql/license-10.sql
/opt/ercole-agent/sql/patch-12.sql
/opt/ercole-agent/sql/edition.sql
/opt/ercole-agent/sql/dbstatus.sql
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
%{_unitdir}/ercole-server.service
%{_presetdir}/60-ercole-server.preset

%changelog
* Mon May  7 2018 Simone Rota <srota2@sorint.it>
- 
