forwarders {
  {{ range service "dns.blocky" }}
  {{ .Address }} port {{ .Port }};
  {{ else }}
  8.8.8.8;
  8.8.4.4;
  {{ end }}
};
