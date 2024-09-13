import {get, patch} from "@/plugins/request"

const clonesetUrl = (cluster_name) => {
    return `/api/v1/proxy/${cluster_name}/k8s/apis/apps.kruise.io/v1alpha1/clonesets`
}
const clonesetWithNsUrl = (cluster_name, namespaces) => {
    return `/api/v1/proxy/${cluster_name}/k8s/apis/apps.kruise.io/v1alpha1/namespaces/${namespaces}/cloneset`
}

export function listClonesets (cluster_name, currentPage, pageSize) {
    let url = clonesetUrl(cluster_name)
    if (currentPage && pageSize) {
        let params = {pageNum: currentPage, pageSize: pageSize }
        return get(url, params)
    }
    return get(url)
}

export function listClonesetsByNs (cluster_name, namespace) {
    return get(`${clonesetWithNsUrl(cluster_name, namespace)}`)
}

export function scaleClonesets (cluster_name, namespace, deployment, data) {
    return patch(`${clonesetWithNsUrl(cluster_name, namespace)}/${deployment}/scale`, data)
}

export function patchClonesets (cluster_name, namespace, deployment, data) {
    return patch(`${clonesetWithNsUrl(cluster_name, namespace)}/${deployment}`, data)
}