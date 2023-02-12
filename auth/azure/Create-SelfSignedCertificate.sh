#!/bin/bash

name=MyCert
exp=1825
pass=MyPassword

read -p "Enter certificate name [ $name ]: " in_name
if [ -n "$in_name" ]; then
  name="$in_name"
fi

read -p "Enter expiration (days) [ $exp ]: " in_exp
if [ -n "$in_exp" ]; then
  exp=$in_exp
fi

read -p "Enter export password [ $pass ]: " in_pass
if [ -n "$in_pass" ]; then
  pass=$in_pass
fi

openssl req -newkey rsa:4096 -nodes -keyout "$name.key" -out "$name.csr" -subj "/CN=localhost"
openssl x509 -signkey "$name.key" -in "$name.csr" -req -days $exp -out "$name.cer"
openssl pkcs12 -export -out "$name.pfx" -inkey "$name.key" -in "$name.cer" -password pass:$pass

echo $pass > $name.txt

rm $name.csr $name.key