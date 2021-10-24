#!/bin/sh

# PROVIDE: livebook
# REQUIRE: DAEMON
# KEYWORD: shutdown
#
# Add the following lines to /etc/rc.conf.local or /etc/rc.conf
# to enable this service:
#
# livebook_enable (bool):	Set it to YES to enable livebook.
#     Default is "NO".
# livebook_user (user):	Set user to run livebook.
#     Default is "livebook".
# livebook_group (group):	Set group to run livebook.
#     Default is "livebook".
# livebook_listen_ip (str): Set IP address for livebook server to listen on.
#     Default is "127.0.0.1".
# livebook_listen_port (int): Set port for livebook server to listen on.
#     Default is "8080".
# livebook_root_data_path (str): Set root data path for livebook server.
#     Default is "/var/db/livebook".
# livebook_syslog_output_enable (bool):	Set to enable syslog output.
#     Default is "NO". See daemon(8).
# livebook_syslog_output_priority (str):	Set syslog priority if syslog enabled.
#     Default is "info". See daemon(8).
# livebook_syslog_output_facility (str):	Set syslog facility if syslog enabled.
#     Default is "daemon". See daemon(8).

. /etc/rc.subr

name=livebook
rcvar=livebook_enable

load_rc_config $name

: ${livebook_enable:="NO"}
: ${livebook_user:="livebook"}
: ${livebook_group:="livebook"}
: ${livebook_listen_ip:="127.0.0.1"}
: ${livebook_listen_port:="8080"}
: ${livebook_root_data_path:="/var/db/livebook"}

DAEMON=$(/usr/sbin/daemon 2>&1 | grep -q syslog ; echo $?)
if [ ${DAEMON} -eq 0 ]; then
        : ${livebook_syslog_output_enable:="NO"}
        : ${livebook_syslog_output_priority:="info"}
        : ${livebook_syslog_output_facility:="daemon"}
        if checkyesno livebook_syslog_output_enable; then
                livebook_syslog_output_flags="-T ${name}"

                if [ -n "${livebook_syslog_output_priority}" ]; then
                        livebook_syslog_output_flags="${livebook_syslog_output_flags} -s ${livebook_syslog_output_priority}"
                fi

                if [ -n "${livebook_syslog_output_facility}" ]; then
                        livebook_syslog_output_flags="${livebook_syslog_output_flags} -l ${livebook_syslog_output_facility}"
                fi
        fi
else
        livebook_syslog_output_enable="NO"
        livebook_syslog_output_flags=""
fi

pidfile=/var/run/livebook.pid
procname="/opt/livebook/.mix/escripts/livebook"
command="/usr/sbin/daemon"
command_args="-r -f -t ${name} ${livebook_syslog_output_flags} -p ${pidfile} -r /usr/bin/env ${livebook_env} ${procname} server --no-token --ip ${livebook_listen_ip} --port ${livebook_listen_port} --root-path ${livebook_root_data_path}"

livebook_chdir="/opt/livebook"
start_precmd=livebook_startprecmd

livebook_startprecmd()
{
        if [ ! -e ${pidfile} ]; then
                install -o ${livebook_user} -g ${livebook_group} /dev/null ${pidfile};
        fi

        if [ ! -e ${livebook_root_data_path} ]; then
                install -d -o ${livebook_user} -g ${livebook_group} ${livebook_root_data_path};
        fi
}

run_rc_command "$1"
