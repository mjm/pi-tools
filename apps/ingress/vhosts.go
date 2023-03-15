package ingress

var virtualHosts = []virtualHost{
	{
		Name: "auth",
		Upstream: upstream{
			Name:        "vault-proxy",
			ServiceName: "vault-proxy",
			ConnectPort: 2220,
		},
		DisableOAuth: true,
		CustomLocationConfig: `
    proxy_set_header Host                    $host;
    proxy_set_header X-Real-IP               $remote_addr;
    proxy_set_header X-Scheme                $scheme;
`,
	},
	{
		Name: "consul",
		Upstream: upstream{
			Name:        "consul",
			ServiceName: "consul",
			ServicePort: 8500,
		},
	},
	{
		Name: "nomad",
		Upstream: upstream{
			Name:        "nomad",
			ServiceName: "http.nomad",
			IPHash:      true,
			Secure:      true,
		},
		CustomLocationConfig: `
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

    proxy_ssl_certificate /etc/nginx/ssl/nomad.pem;
    proxy_ssl_certificate_key /etc/nginx/ssl/nomad.pem;
    proxy_ssl_trusted_certificate /etc/nginx/ssl/nomad.ca.crt;

    # Nomad blocking queries will remain open for a default of 5 minutes.
    # Increase the proxy timeout to accommodate this timeout with an
    # additional grace period.
    proxy_read_timeout 310s;

    # Nomad log streaming uses streaming HTTP requests. In order to
    # synchronously stream logs from Nomad to NGINX to the browser
    # proxy buffering needs to be turned off.
    proxy_buffering off;

    # The Upgrade and Connection headers are used to establish
    # a WebSockets connection.
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";

    # The default Origin header will be the proxy address, which
    # will be rejected by Nomad. It must be rewritten to be the
    # host address instead.
    proxy_set_header Origin "${scheme}://${proxy_host}";
`,
	},
	{
		Name: "vault",
		Upstream: upstream{
			Name:        "vault",
			ServiceName: "vault",
		},
	},
	{
		Name: "prometheus",
		Upstream: upstream{
			Name:        "prometheus",
			ServiceName: "prometheus",
		},
	},
	{
		Name: "alertmanager",
		Upstream: upstream{
			Name:        "alertmanager",
			ServiceName: "alertmanager",
		},
	},
	{
		Name: "grafana",
		Upstream: upstream{
			Name:        "grafana",
			ServiceName: "grafana",
			ConnectPort: 3000,
		},
	},
	{
		Name: "go",
		Upstream: upstream{
			Name:        "go-links",
			ServiceName: "go-links",
			ConnectPort: 4240,
		},
		CustomServerConfig: `
  add_header Strict-Transport-Security "max-age=2628000" always;
`,
	},
	{
		Name: "minio",
		Upstream: upstream{
			Name:        "minio",
			ServiceName: "minio",
		},
		DisableOAuth: true,
		CustomServerConfig: `
  # To allow special characters in headers
  ignore_invalid_headers off;
  # Allow any size file to be uploaded.
  # Set to a value such as 1000m; to restrict file size to a specific value
  client_max_body_size 0;
  # To disable buffering
  proxy_buffering off;
`,
		CustomLocationConfig: `
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header Host $http_host;

    proxy_connect_timeout 300;
    # Default is HTTP/1, keepalive is only enabled in HTTP/1.1
    proxy_http_version 1.1;
    proxy_set_header Connection "";
    chunked_transfer_encoding off;
`,
	},
	{
		Name: "minio-console",
		Upstream: upstream{
			Name:        "minio-console",
			ServiceName: "minio-console",
		},
		DisableOAuth: true,
		CustomServerConfig: `
  # To allow special characters in headers
  ignore_invalid_headers off;
  # Allow any size file to be uploaded.
  # Set to a value such as 1000m; to restrict file size to a specific value
  client_max_body_size 0;
  # To disable buffering
  proxy_buffering off;
`,
		CustomLocationConfig: `
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header Host $http_host;

    proxy_connect_timeout 300;
    # Default is HTTP/1, keepalive is only enabled in HTTP/1.1
    proxy_http_version 1.1;
    proxy_set_header Connection "";
    chunked_transfer_encoding off;
`,
	},
	{
		Name: "pkg",
		Upstream: upstream{
			Name:        "poudriere-web",
			ServiceName: "poudriere-web",
		},
		DisableOAuth: true,
	},
	{
		Name: "paperless",
		Upstream: upstream{
			Name:        "paperless",
			ServiceName: "paperless",
		},
		CustomServerConfig: `
  client_max_body_size 100M;
`,
		CustomLocationConfig: `
    # These configuration options are required for WebSockets to work.
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";

    proxy_redirect off;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Host $server_name;
`,
	},
	{
		Name: "code",
		Upstream: upstream{
			Name:        "phabricator",
			ServiceName: "phabricator",
		},
		DisableOAuth: true,
		CustomServerConfig: `
  # To allow special characters in headers
  ignore_invalid_headers off;
  # Allow any size file to be uploaded.
  # Set to a value such as 1000m; to restrict file size to a specific value
  client_max_body_size 0;
  # To disable buffering
  proxy_buffering off;
  proxy_connect_timeout       300;
  proxy_send_timeout          300;
  proxy_read_timeout          300;
  send_timeout                300;
`,
		CustomLocationConfig: `
    proxy_redirect off;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Host $server_name;
`,
	},
	{
		Name: "ci",
		Upstream: upstream{
			Name:        "teamcity",
			ServiceName: "teamcity",
		},
		DisableOAuth: true,
		CustomServerConfig: `
  proxy_read_timeout     1200;
  proxy_connect_timeout  240;
  client_max_body_size   0;    # maximum size of an HTTP request. 0 allows uploading large artifacts to TeamCity
`,
		CustomLocationConfig: `
    proxy_http_version  1.1;
    proxy_set_header    Host $server_name:$server_port;
    proxy_set_header    X-Forwarded-Host $http_host;    # necessary for proper absolute redirects and TeamCity CSRF check
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    X-Forwarded-For $remote_addr;
    proxy_set_header    Upgrade $http_upgrade; # WebSocket support
    proxy_set_header    Connection $connection_upgrade; # WebSocket support
`,
	},
	{
		Name: "homelab",
		Upstream: upstream{
			Name:        "homelab",
			ServiceName: "homelab",
		},
		// skip auth redirect for /app endpoints
		CustomServerConfig: `
  location /app {
    proxy_pass http://homelab;
    proxy_http_version  1.1;
    proxy_set_header    Host $server_name:$server_port;
    proxy_set_header    X-Forwarded-Host $http_host;
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    X-Forwarded-For $remote_addr;
  }
`,
		CustomLocationConfig: `
    proxy_http_version  1.1;
    proxy_set_header    Host $server_name:$server_port;
    proxy_set_header    X-Forwarded-Host $http_host;
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    X-Forwarded-For $remote_addr;
    proxy_set_header    Upgrade $http_upgrade; # WebSocket support
    proxy_set_header    Connection $connection_upgrade; # WebSocket support
`,
	},
	{
		Name: "livebook",
		Upstream: upstream{
			Name:        "livebook",
			ServiceName: "livebook",
		},
		CustomLocationConfig: `
    proxy_http_version  1.1;
    proxy_set_header    Host $server_name:$server_port;
    proxy_set_header    X-Forwarded-Host $http_host;
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    X-Forwarded-For $remote_addr;
    proxy_set_header    Upgrade $http_upgrade; # WebSocket support
    proxy_set_header    Connection $connection_upgrade; # WebSocket support
`,
	},
	{
		Name: "adminer",
		Upstream: upstream{
			Name:        "adminer",
			ServiceName: "adminer",
			ConnectPort: 10000,
		},
	},
	{
		Name: "guacamole",
		Upstream: upstream{
			Name:        "guacamole",
			Path:        "/guacamole/",
			ServiceName: "guacamole",
		},
		CustomServerConfig: `
  # Allow any size file to be uploaded.
  # Set to a value such as 1000m; to restrict file size to a specific value
  client_max_body_size 0;
  # To disable buffering
  proxy_buffering off;
`,
		CustomLocationConfig: `
    proxy_http_version  1.1;
    proxy_set_header    Host $server_name:$server_port;
    proxy_set_header    X-Forwarded-Host $http_host;
    proxy_set_header    X-Forwarded-Proto $scheme;
    proxy_set_header    X-Forwarded-For $remote_addr;
    proxy_set_header    Upgrade $http_upgrade; # WebSocket support
    proxy_set_header    Connection $connection_upgrade; # WebSocket support
`,
	},
}

type virtualHost struct {
	Name                 string
	Upstream             upstream
	DisableOAuth         bool
	CustomServerConfig   string
	CustomLocationConfig string
}
