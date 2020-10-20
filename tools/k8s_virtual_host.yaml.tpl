apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: {NAMESPACE}
  name: {NAME}
spec:
  rules:
    - host: {NAME}
      http:
        paths:
          - path: "/"
            pathType: Prefix
            backend:
              service:
                name: {SERVICE_NAME}
                port:
                  number: {PORT}
