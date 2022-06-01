locals {
  subnet_settings = {
    offset = 8 # 8 bit offset

    # This bumps the subnet CIDRS to space them accordingly
    public = 0
    database = 16
    private = 32
    intra = 64
  }
  default_tags = {
    "cost-center" = "psychic-guacamole"
    "environment" = "development"
    "managed-by"  = "terraform"
  }
}
