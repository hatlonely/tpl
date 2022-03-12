package rpcx

var opsYaml = `
name: {{ .Name }}

dep:
  ops:
    type: git
    url: "https://github.com/hatlonely/ops.git"
    version: master

env:
  default:
    NAME: "{{ .Name }}"
	GIT_TAG: "$(git describe --tags)"
    VERSION: "$(git describe --tags | awk '{print(substr($0,2,length($0)))}')"
    REGISTRY_ENDPOINT: "{{ "{{ .registry.endpoint }}" }}"
    REGISTRY_USERNAME: "{{ "{{ .registry.username }}" }}"
    REGISTRY_PASSWORD: "{{ "{{ .registry.password }}" }}"
    REGISTRY_NAMESPACE: "{{ "{{ .registry.namespace }}" }}"
    {{ if not .Ops.EnableHelm }}
    K8S_CONTEXT: "home-k8s"
    NAMESPACE: "dev"
    PULL_SECRET_NAME: "hatlonely-pull-secret"
    REPLICA_COUNT: 2
    INGRESS_HOST: "k8s.rpc.tool.hatlonely.com"
    SECRET_NAME: "rpc-tool-tls"
    {{ end }}


task:
  image:
    step:
      - make image
      - docker login --username="${REGISTRY_USERNAME}" --password="${REGISTRY_PASSWORD}" "${REGISTRY_ENDPOINT}"
      - docker push "${REGISTRY_ENDPOINT}/${REGISTRY_NAMESPACE}/${NAME}:${VERSION}"
  {{ if not .Ops.EnableHelm }}
  helm:
    args:
      cmd:
        type: string
        default: diff
        validation: x in ["diff", "install", "upgrade", "delete"]
    step:
      - test "${K8S_CONTEXT}" == "$(kubectl config current-context)" || exit 1
      - sh ${DEP}/ops/tool/render.sh ${DEP}/ops/helm/rpc-app ${TMP}/helm/${NAME}
      - sh ${DEP}/ops/tool/render.sh ops/helm/values-adapter.yaml.tpl ${TMP}/helm/${NAME}/values-adapter.yaml
      - |
        case "${cmd}" in
          "diff"|"") helm diff upgrade "${NAME}" -n "${NAMESPACE}" "${TMP}/helm/${NAME}" -f "${TMP}/helm/${NAME}/values-adapter.yaml" --allow-unreleased;;
          "install") helm install "${NAME}" -n "${NAMESPACE}" "${TMP}/helm/${NAME}" -f "${TMP}/helm/${NAME}/values-adapter.yaml";;
          "upgrade") helm upgrade "${NAME}" -n "${NAMESPACE}" "${TMP}/helm/${NAME}" -f "${TMP}/helm/${NAME}/values-adapter.yaml";;
          "delete") helm delete "${NAME}" -n "${NAMESPACE}";;
        esac
  {{ end }}
`
