forwarders {
  {{ range service "dns.pihole" }}
  {{ .Address }} port {{ .Port }};
  {{ else }}
  8.8.8.8;
  8.8.4.4;
  {{ end }}
};
