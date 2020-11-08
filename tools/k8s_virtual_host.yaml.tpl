apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: {NAMESPACE}
  name: {NAME}
  annotations:
    nginx.ingress.kubernetes.io/auth-url: "http://oauth2-proxy.auth.svc.cluster.local/oauth2/auth"
    nginx.ingress.kubernetes.io/auth-response-headers: "X-Auth-Request-User,X-Auth-Request-Email"
    nginx.ingress.kubernetes.io/auth-signin: "http://homebase.homelab/oauth2/start?rd=$escaped_request_uri"
    cert-manager.io/cluster-issuer: ca-issuer
spec:
  rules:
    - host: {NAME}.homelab
      http:
        paths:
          - path: "/"
            pathType: Prefix
            backend:
              service:
                name: {SERVICE_NAME}
                port:
                  number: {PORT}
  tls:
    - hosts:
        - {NAME}.homelab
      secretName: {NAME}-cert
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: {NAMESPACE}
  name: {NAME}-redirect
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: "http://{NAME}.homelab/$1"
    cert-manager.io/cluster-issuer: ca-issuer
spec:
  rules:
    - host: {NAME}
      http:
        paths:
            - path: "/(.*)"
              pathType: Prefix
              backend:
                service:
                  name: {SERVICE_NAME}
                  port:
                    number: {PORT}
  tls:
    - hosts:
        - {NAME}
      secretName: {NAME}-redirect-cert
