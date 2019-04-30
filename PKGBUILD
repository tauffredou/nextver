# Maintainer: Thomas Auffredou <thomas.auffredou@gmail.com>
# Contributor: Adrien Folie <folie.adrien@gmail.com>
_pkgname=nextver
pkgname="${_pkgname}-git"
pkgver=0.0.0
pkgrel=1
pkgdesc="Calculates the next version from the git history, using the \"conventional commits\" specification"
arch=("x86_64")
url="https://github.com/tauffredou/nextver"
makedepends=("go")
source=("${_pkgname}::git+${url}.git")
md5sums=('SKIP')

build() {
	cd "$srcdir/$_pkgname"
	go mod download
	go build
}

package() {
	cd "$srcdir/$_pkgname"

	install -Dm755 "$_pkgname" "$pkgdir/usr/bin/$_pkgname"
}
