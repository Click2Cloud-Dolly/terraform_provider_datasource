terraform {
  required_providers {
    employee = {
      version = "0.2"
      source  = "hashicorp.com/edu/employee"
    }
  }
}

data "hashicups_employee" "employees" {
  provider = employee
}

# Returns all coffees
output "all_employee" {
  value = data.hashicups_employee.employees.*

}
