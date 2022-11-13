#!/bin/sh

# PROVIDE: teamcity-agent
# REQUIRE: DAEMON
# KEYWORD: shutdown
#
# Add the following lines to /etc/rc.conf.local or /etc/rc.conf
# to enable this service:
#
# teamcity_agent_enable (bool):	Set it to YES to enable TeamCity build agent.
#     Default is "NO".
# teamcity_agent_user (user):	Set user to run TeamCity agent.
#     Default is "teamcity".
# teamcity_agent_group (group):	Set group to run TeamCity agent.
#     Default is "teamcity".
# teamcity_agent_jre_home (path): Set path to JRE home for Java installation.
#     Default is "/usr/local/openjdk11-jre".
# teamcity_agent_syslog_output_enable (bool):	Set to enable syslog output.
#     Default is "NO". See daemon(8).
# teamcity_agent_syslog_output_priority (str):	Set syslog priority if syslog enabled.
#     Default is "info". See daemon(8).
# teamcity_agent_syslog_output_facility (str):	Set syslog facility if syslog enabled.
#     Default is "daemon". See daemon(8).

. /etc/rc.subr

name=teamcity_agent
rcvar=teamcity_agent_enable

load_rc_config $name

: ${teamcity_agent_enable:="NO"}
: ${teamcity_agent_user:="teamcity"}
: ${teamcity_agent_group:="teamcity"}
: ${teamcity_agent_jre_home:="/usr/local/openjdk11-jre"}

DAEMON=$(/usr/sbin/daemon 2>&1 | grep -q syslog ; echo $?)
if [ ${DAEMON} -eq 0 ]; then
        : ${teamcity_agent_syslog_output_enable:="NO"}
        : ${teamcity_agent_syslog_output_priority:="info"}
        : ${teamcity_agent_syslog_output_facility:="daemon"}
        if checkyesno teamcity_agent_syslog_output_enable; then
                teamcity_agent_syslog_output_flags="-T ${name}"

                if [ -n "${teamcity_agent_syslog_output_priority}" ]; then
                        teamcity_agent_syslog_output_flags="${teamcity_agent_syslog_output_flags} -s ${teamcity_agent_syslog_output_priority}"
                fi

                if [ -n "${teamcity_agent_syslog_output_facility}" ]; then
                        teamcity_agent_syslog_output_flags="${teamcity_agent_syslog_output_flags} -l ${teamcity_agent_syslog_output_facility}"
                fi
        fi
else
        teamcity_agent_syslog_output_enable="NO"
        teamcity_agent_syslog_output_flags=""
fi

procname="bin/agent.sh"
#command="/usr/sbin/daemon"
#command_args="-f -t ${name} ${teamcity_syslog_output_flags} -p ${pidfile} -r /usr/bin/env HOME=/opt/TeamCity USER=${teamcity_user} TEAMCITY_DATA_PATH=${teamcity_data_path} JRE_HOME=${teamcity_jre_home} ${teamcity_env} ${procname} run"

start_cmd=teamcity_agent_start
stop_cmd=teamcity_agent_stop

teamcity_agent_start()
{
        cd /opt/TeamCity/buildAgent
        su -m teamcity -c "/usr/bin/env HOME=/opt/TeamCity USER=${teamcity_agent_user} JRE_HOME=${teamcity_agent_jre_home} ${teamcity_agent_env} /usr/local/bin/bash ${procname} start"
}

teamcity_agent_stop()
{
        cd /opt/TeamCity/buildAgent
        su -m teamcity -c "/usr/bin/env HOME=/opt/TeamCity USER=${teamcity_agent_user} JRE_HOME=${teamcity_agent_jre_home} ${teamcity_agent_env} /usr/local/bin/bash ${procname} stop"
}

run_rc_command "$1"
