provider "kubernetes" {
  config_path = "~/.kube/config"
}

resource "kubernetes_service" "node_exporter" {
  metadata {
    name = "maxscale-exporter"
  }
  spec {
    selector = {
      app = kubernetes_pod.node_exporter.metadata.0.labels.app
    }

    port {
      port = 9104
      name = "metrics"
    }

    type = "ClusterIP"
  }
}

resource "kubernetes_pod" "node_exporter" {
  metadata {
    name = "maxscale-exporter"
    labels = {
      app = "maxscale-exporter"
    }
  }

  spec {
    container {
      image             = "maxscale_exporter"
      name              = "maxscale-exporter"
      image_pull_policy = "IfNotPresent"

      port {
        port = 9104
        name = "metrics"
      }
    }
  }
}
