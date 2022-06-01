resource "aws_db_instance" "psychic_guacamole" {
  db_name              = "psychicguacamole"

  engine               = "mysql"
  engine_version       = "8.0"
  instance_class       = "db.t4g.micro"

  username             = "foo"
  password             = "foobarbaz"

  parameter_group_name = "default.mysql8.0"
  skip_final_snapshot   = true

  publicly_accessible = true

  # Storage
  allocated_storage    = 5
  storage_type         = "standard"

  vpc_security_group_ids = [aws_security_group.rds.id]
  db_subnet_group_name   = module.vpc.database_subnet_group_name

  lifecycle {
    ignore_changes = [username, password]
  }
}
