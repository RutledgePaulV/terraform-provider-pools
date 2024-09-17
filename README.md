# Terraform Pools Provider

A terraform provider offering stateful resource pools. A pool contains
resources and allocates the resources therein to borrowers. Borrowers
consistently receive the same resource from the pool until either the
borrower or the resource are removed from the pool.

This is useful for random exclusive allocation of a finite set of resources.
