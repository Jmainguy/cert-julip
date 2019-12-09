# cert-julip
## Purpose
To provide a bridge between cert-manager and openshift
## History
cert-manager is great at generating, and managing certificates, it stores these certificates in secrets, which kubernetes ingress can then use.

Openshift however does not allow secrets in their routes, so it is unable to effectively use cert-manager in an automated fashion.

cert-julip looks for the certificate: label in a route, where the value will be the cert-manager kind=certificate you want to use.

cert-julip gets the secret from the certificate, and then adds the certs and key from the secret into the route.

### How To
1. Setup a certificate with cert-manager first. Yaml will look similiar to
```
---
apiVersion: certmanager.k8s.io/v1alpha1
kind: Certificate
metadata:
  name: jmainguy-example-com
  namespace: jmainguy
spec:
  secretName: jmainguy-example-com-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  commonName: 'jmainguy.example.com'
  dnsNames:
  - jmainguy.example.com
  acme:
    config:
    - dns01:
        provider: route53
      domains:
      - jmainguy.example.com
```
2. Edit openshift route and add label
```
apiVersion: route.openshift.io/v1
kind: Route
metadata:
  labels:
    certificate: jmainguy-example-com
```
3. cert-julip will auto-populate the certificate, key, and ca certificate from the cert-manager certificate you linked to above. This works with edge, and reencrypt routes
