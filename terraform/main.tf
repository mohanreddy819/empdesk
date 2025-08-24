resource "google_compute_instance" "go_app_vm" {
  name         = "go-app-vm"
  machine_type = "e2-micro"       
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2004-lts"
    }
  }

  network_interface {
    network       = "default"
    access_config {}
  }
}