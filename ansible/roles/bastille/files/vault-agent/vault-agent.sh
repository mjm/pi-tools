#!/bin/sh

# PROVIDE: vault-agent
# REQUIRE: DAEMON
# KEYWORD: shutdown
#
# Add the following lines to /etc/rc.conf.local or /etc/rc.conf
# to enable this service:
#
# vault_agent_enable (bool):	Set it to YES to enable vault.
#			Default is "NO".
# vault_agent_user (user):	Set user to run vault.
#			Default is "vault".
# vault_agent_group (group):	Set group to run vault.
#			Default is "vault".
# vault_agent_config (file):	Set vault config file.
#			Default is "/usr/local/etc/vault-agent.hcl".
# vault_agent_syslog_output_enable (bool):	Set to enable syslog output.
#					Default is "NO". See daemon(8).
# vault_agent_syslog_output_priority (str):	Set syslog priority if syslog enabled.
#					Default is "info". See daemon(8).
# vault_agent_syslog_output_facility (str):	Set syslog facility if syslog enabled.
#					Default is "daemon". See daemon(8).

. /etc/rc.subr

name=vault_agent
rcvar=vault_agent_enable

load_rc_config $name

: ${vault_agent_enable:="NO"}
: ${vault_agent_user:="vault"}
: ${vault_agent_group:="vault"}
: ${vault_agent_config:="/usr/local/etc/vault-agent.hcl"}

DAEMON=$(/usr/sbin/daemon 2>&1 | grep -q syslog ; echo $?)
if [ ${DAEMON} -eq 0 ]; then
        : ${vault_agent_syslog_output_enable:="NO"}
        : ${vault_agent_syslog_output_priority:="info"}
        : ${vault_agent_syslog_output_facility:="daemon"}
        if checkyesno vault_agent_syslog_output_enable; then
                vault_agent_syslog_output_flags="-T ${name}"

                if [ -n "${vault_agent_syslog_output_priority}" ]; then
                        vault_agent_syslog_output_flags="${vault_agent_syslog_output_flags} -s ${vault_agent_syslog_output_priority}"
                fi

                if [ -n "${vault_agent_syslog_output_facility}" ]; then
                        vault_agent_syslog_output_flags="${vault_agent_syslog_output_flags} -l ${vault_agent_syslog_output_facility}"
                fi
        fi
else
        vault_agent_syslog_output_enable="NO"
        vault_agent_syslog_output_flags=""
fi

pidfile=/var/run/vault-agent.pid
procname="/usr/local/bin/vault"
command="/usr/sbin/daemon"
command_args="-f -t ${name} ${vault_agent_syslog_output_flags} -p ${pidfile} -r /usr/bin/env ${vault_agent_env} ${procname} agent -config=${vault_agent_config}"

extra_commands="reload monitor"
monitor_cmd=vault_agent_monitor
start_precmd=vault_agent_startprecmd
required_files="$vault_agent_config"

vault_agent_monitor()
{
	sig_reload=USR1
	run_rc_command "reload"
}

vault_agent_startprecmd()
{
        if [ ! -e ${pidfile} ]; then
                install -o ${vault_agent_user} -g ${vault_agent_group} /dev/null ${pidfile};
        fi
}

run_rc_command "$1"
