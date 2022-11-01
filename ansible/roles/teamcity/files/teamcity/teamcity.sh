#!/bin/sh

# PROVIDE: teamcity
# REQUIRE: DAEMON
# KEYWORD: shutdown
#
# Add the following lines to /etc/rc.conf.local or /etc/rc.conf
# to enable this service:
#
# teamcity_enable (bool):	Set it to YES to enable TeamCity.
#     Default is "NO".
# teamcity_user (user):	Set user to run TeamCity.
#     Default is "teamcity".
# teamcity_group (group):	Set group to run TeamCity.
#     Default is "teamcity".
# teamcity_data_path (path): Set path to store TeamCity data.
#     Default is "/var/db/teamcity".
# teamcity_jre_home (path): Set path to JRE home for Java installation.
#     Default is "/usr/local/openjdk11-jre".
# teamcity_syslog_output_enable (bool):	Set to enable syslog output.
#     Default is "NO". See daemon(8).
# teamcity_syslog_output_priority (str):	Set syslog priority if syslog enabled.
#     Default is "info". See daemon(8).
# teamcity_syslog_output_facility (str):	Set syslog facility if syslog enabled.
#     Default is "daemon". See daemon(8).

. /etc/rc.subr

name=teamcity
rcvar=teamcity_enable

load_rc_config $name

: ${teamcity_enable:="NO"}
: ${teamcity_user:="teamcity"}
: ${teamcity_group:="teamcity"}
: ${teamcity_data_path:="/var/db/teamcity"}
: ${teamcity_jre_home:="/usr/local/openjdk11-jre"}

DAEMON=$(/usr/sbin/daemon 2>&1 | grep -q syslog ; echo $?)
if [ ${DAEMON} -eq 0 ]; then
        : ${teamcity_syslog_output_enable:="NO"}
        : ${teamcity_syslog_output_priority:="info"}
        : ${teamcity_syslog_output_facility:="daemon"}
        if checkyesno teamcity_syslog_output_enable; then
                teamcity_syslog_output_flags="-T ${name}"

                if [ -n "${teamcity_syslog_output_priority}" ]; then
                        teamcity_syslog_output_flags="${teamcity_syslog_output_flags} -s ${teamcity_syslog_output_priority}"
                fi

                if [ -n "${teamcity_syslog_output_facility}" ]; then
                        teamcity_syslog_output_flags="${teamcity_syslog_output_flags} -l ${teamcity_syslog_output_facility}"
                fi
        fi
else
        teamcity_syslog_output_enable="NO"
        teamcity_syslog_output_flags=""
fi

procname="bin/teamcity-server.sh"

start_precmd=teamcity_startprecmd
start_cmd=teamcity_start
stop_cmd=teamcity_stop

teamcity_startprecmd()
{
        if [ ! -e ${teamcity_data_path} ]; then
                install -d -o ${teamcity_user} -g ${teamcity_group} ${teamcity_data_path};
        fi
}

teamcity_start()
{
        cd /opt/TeamCity
        /usr/bin/env HOME=/opt/TeamCity USER=${teamcity_user} TEAMCITY_DATA_PATH=${teamcity_data_path} JRE_HOME=${teamcity_jre_home} ${teamcity_env} ${procname} start
}

teamcity_stop()
{
        cd /opt/TeamCity
        /usr/bin/env HOME=/opt/TeamCity USER=${teamcity_user} TEAMCITY_DATA_PATH=${teamcity_data_path} JRE_HOME=${teamcity_jre_home} ${teamcity_env} ${procname} stop
}

run_rc_command "$1"
