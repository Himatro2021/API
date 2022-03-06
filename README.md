# API CONTRACT

1. /login
  - POST
  - Content-Type: application/x-www-form-urlencoded
  - Form:
    <br>a. NPM (string)
    <br>b. password (string)
  - Returns: token

2. /admin
  - GET
  - Require login
  - Authentication method: Bearer

3. /admin/absensi
  - POST
  - Content-Type: application/x-www-form-urlencoded
  - Form:
  <br>a. title (string)
  <br>b. startAtDate (string) -> must be in format yyyy-mm-dd
  <br>c. startAtTime (string) -> must be in format hh:mm:ss
  <br>d. finishAtDate (string) -> must be in format yyyy-mm-dd
  <br>e. finishAtTime (string) -> must be in format hh:mm:ss
  <br>f. participant (string) -> must be one of defined departement name listend below. Case insensitive

4. /absensi/:absentID
  - GET
  - Params:
  <br>a. absentID (number)
  - Returns: absent result -> {npm. updatedAt, keterangan, nama, departemen }

## Host
api.himatro.luckyakbar.tech

## Defined Departement Name
1. Pengurus Harian -> PH
2. Departemen Pendidikan dan Pengembangan Diri -> PPD
3. Departemen Kaderisasi Dan Pengembangan Organisasi -> KPO
4. Departemen Komunikasi Dan Informasi -> KOMINFO
5. Departemen Sosial dan Kewirausahaan -> KWU
6. Departemen Pengembangan Keteknikan -> BANGTEK