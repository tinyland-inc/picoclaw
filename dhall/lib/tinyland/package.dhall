-- Tinyland shared Dhall library.
--
-- Common types and fragments shared across the Tinyland ecosystem:
-- tinyclaw, tailnet-acl, remote-juggler, gloriousflywheel.

let ACL = ./acl.dhall
let AperturePolicy = ./aperture-policy.dhall

in  { ACL, AperturePolicy }
