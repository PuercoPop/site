# 2023-08-25

Found about guix deploy, that takes care of provisioning the VM (droplet) with
an 'infected' guix system. Infected means that it uses another distro, debian,
as the initial distro and then installs guix-sd there.

https://guix.gnu.org/manual/en/html_node/Invoking-guix-deploy.html
https://guix.gnu.org/fr/blog/2019/managing-servers-with-gnu-guix-a-tutorial/
https://github.com/clojure-quant/infra-guix/blob/11c547e7b97f92ce41064e0333bc101abe2d50cd/os/ocean-deploy.scm#L1
https://stumbles.id.au/getting-started-with-guix-deploy.html
