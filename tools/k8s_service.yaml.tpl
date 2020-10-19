apiVersion: v1
kind: Service
metadata:
  namespace: {NAMESPACE}
  name: {NAME}
spec:
  selector:
    app: {APP}
  ports:
    - port: 80
      targetPort: {TARGET_PORT}
