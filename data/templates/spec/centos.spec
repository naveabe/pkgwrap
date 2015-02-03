#
# Spec for Package {{.Name}} {{.Version}}
#
# Copyright (c) {{.Year}} {{.Packager}}.
#
# License {{.License}} 
#

# Turn off python byte compiling
%global __os_install_post %(echo '%{__os_install_post}' | sed -e 's!/usr/lib[^[:space:]]*/brp-python-bytecompile[[:space:]].*$!!g')

%define NAME     {{.Name}}
%define VERSION  {{.Version}}
%define RELEASE  {{.Release}}
%define PACKAGER {{.Packager}}

Name            : %{NAME}
Version         : %{VERSION}
Release         : %{RELEASE}%{?dist}
Packager        : %{PACKAGER}{{if .Url}}
Url             : {{.Url}}{{end}}
Summary         : {{.Summary}}
License         : {{.License}}
Group           : {{.Group}}
Source0         : {{.Source}}{{if .BuildRequires}}
BuildRequires   : {{.BuildRequires}}{{end}}{{if .Requires}}
Requires        : {{.Requires}}{{end}}

%description
{{.Description}}

{{if .Prep}}%prep
{{range .Prep}}{{.}}
{{end}}{{end}}

{{if .Build}}%build
{{range .Build}}{{.}}
{{end}}{{end}}

{{if .PreInstall}}%pre
{{range .PreInstall}}{{.}}
{{end}}{{end}}

%install
{{range .Install}}{{.}}
{{end}}

{{if .PostInstall}}%post
{{range .PostInstall}}{{.}}
{{end}}{{end}}

{{if .PreUninstall}}%preun
{{range .PreUninstall}}{{.}}
{{end}}{{end}}

{{if .PostUninstall}}%postun
{{range .PostUninstall}}{{.}}
{{end}}{{end}}

%clean
{{range .Clean}}{{.}}
{{end}}

%files
{{range .Files}}{{.}}
{{end}}