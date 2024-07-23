resource "aria_property_group" "vm_common" {
  name        = "VM_Common"
  description = "Common Machines properties."
  type        = "INPUT"

  properties = {
    vlan = {
      name               = "vlan"
      title              = "VLAN"
      description        = "VLAN"
      type               = "string"
      encrypted          = false
      read_only          = false
      recreate_on_update = false
      one_of = [
        {
          title     = "734b"
          const     = "734b,gold"
          encrypted = false
        },
        {
          title     = "755"
          const     = "755,gold"
          encrypted = false
        }
      ]
    }
    srv_role = {
      name               = "srv_role"
      title              = "Rôle du serveur"
      description        = "Rôle du serveur dans l'architecture sur <b>3</b> caractères."
      type               = "string"
      encrypted          = false
      read_only          = false
      recreate_on_update = false
      min_length         = 3
      max_length         = 3
      one_of = [
        {
          title     = "Front"
          const     = "FRT"
          encrypted = false
        },
        {
          title     = "Backend"
          const     = "BKD"
          encrypted = false
        },
        {
          title     = "Base De Données"
          const     = "BDD"
          encrypted = false
        },
        {
          title     = "Message Broker"
          const     = "BRK"
          encrypted = false
        },
        {
          title     = "Monitoring"
          const     = "MON"
          encrypted = false
        },
        {
          title     = "Reverse Proxy"
          const     = "RPX"
          encrypted = false
        },
        {
          title     = "Cache"
          const     = "CCH"
          encrypted = false
        },
        {
          title     = "Générique"
          const     = "GEN"
          encrypted = false
        }
      ]
    }
  }
}
