apiVersion: v1
kind: Service
metadata:
  namespace: {NAMESPACE}
  name: {NAME}
  labels:
    app: {APP}
spec:
  selector:
    app: {APP}
  ports:
    - name: http
      port: 80
      targetPort: {TARGET_PORT}
      protocol: TCP
