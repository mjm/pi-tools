<?php

class AdminerCustom {
    var $servers;

    function __construct($servers) {
        $this->servers = $servers;
        if ($_POST["auth"]) {
            $key = $_POST["auth"]["server"];
            $_POST["auth"]["driver"] = $this->servers[$key]["driver"];
            $_POST["auth"]["db"] = $this->servers[$key]["db"];
        }
    }

    function credentials() {
        $config = $this->servers[SERVER];
        return array($config["server"], $config["username"], $config["password"]);
    }

    function login($login, $password) {
        return true;
    }

	function loginFormField($name, $heading, $value) {
		if ($name == 'driver' || $name == 'username' || $name == 'password' || $name == 'db') {
			return '';
		} elseif ($name == 'server') {
			return $heading . "<select name='auth[server]'>" . optionlist(array_keys($this->servers), SERVER) . "</select>\n";
		}
	}
}

return new AdminerCustom(array(
{{ with secret "database/creds/go-links" }}
    'go-links' => array(
        'driver' => 'pgsql',
        'server' => 'postgresql.service.consul',
        'username' => {{ .Data.username | toJSON }},
        'password' => {{ .Data.password | toJSON }},
        'db' => 'go_links',
    ),
{{ end }}
{{ with secret "database/creds/grafana" }}
    'grafana' => array(
        'driver' => 'pgsql',
        'server' => 'postgresql.service.consul',
        'username' => {{ .Data.username | toJSON }},
        'password' => {{ .Data.password | toJSON }},
        'db' => 'grafana',
    ),
{{ end }}
{{ with secret "database/creds/paperless" }}
    'paperless' => array(
        'driver' => 'pgsql',
        'server' => 'postgresql.service.consul',
        'username' => {{ .Data.username | toJSON }},
        'password' => {{ .Data.password | toJSON }},
        'db' => 'paperless',
    ),
{{ end }}
{{ with secret "database/creds/presence" }}
    'presence' => array(
        'driver' => 'pgsql',
        'server' => 'postgresql.service.consul',
        'username' => {{ .Data.username | toJSON }},
        'password' => {{ .Data.password | toJSON }},
        'db' => 'trips',
    ),
{{ end }}
{{ with secret "database/creds/phabricator" }}
    'phabricator' => array(
        'driver' => 'server',
        'server' => 'mysql.service.consul',
        'username' => {{ .Data.username | toJSON }},
        'password' => {{ .Data.password | toJSON }},
    ),
{{ end }}
));
