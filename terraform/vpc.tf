module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = local.default_tags["environment"]
  cidr = var.vpc_cidr

  azs = data.aws_availability_zones.current_region.names

  igw_tags = local.default_tags

  # This assumes we're creating subnets 8 bits smaller (from /16 to /24 for instance)

  public_subnets = [
    # You're going to end up having /24 networks in a /16 CIDR range = 256 networks

    # cidrsubnet("10.0.0.0/16", 8, 0)
    # 10.0.0.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.public),
    # cidrsubnet("10.0.0.0/16", 8, 1)
    # 10.0.1.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.public + 1),
    # cidrsubnet("10.0.0.0/16", 8, 2)
    # 10.0.2.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.public + 2)
  ]

  public_subnet_tags = merge(
    local.default_tags,
  )

  database_subnets = [

    # cidrsubnet("10.0.0.0/16", 8, 16)
    # 10.0.0.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.database),
    # cidrsubnet("10.0.0.0/16", 8, 17)
    # 10.0.1.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.database + 1),
    # cidrsubnet("10.0.0.0/16", 8, 18)
    # 10.0.2.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.database + 2)
  ]

  database_subnet_tags = merge(
    local.default_tags,
  )

  private_subnets = [
    # cidrsubnet("10.0.0.0/16", 8, 32)
    # 10.0.32.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.private),
    # cidrsubnet("10.0.0.0/16", 8, 33)
    # 10.0.33.0/24
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.private + 1),
    cidrsubnet(var.vpc_cidr, local.subnet_settings.offset, local.subnet_settings.private + 2)
  ]

  private_subnet_tags = merge(
    local.default_tags,
  )

  vpc_tags = merge(
    local.default_tags,
  )

  enable_dns_hostnames = true
  enable_dns_support = true

  enable_nat_gateway     = true
  single_nat_gateway     = false
  one_nat_gateway_per_az = false


  create_database_subnet_group           = true
  create_database_subnet_route_table     = true
  create_database_internet_gateway_route = true

  enable_vpn_gateway = false

  tags = local.default_tags
}
