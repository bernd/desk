#!/usr/bin/fish
printf "\


" | desk init

cp examples/* ~/.desk/desks/

set LIST (desk ls)
echo $LIST | grep "hello" >/dev/null
test $status -ne 0; and echo "hello desk not found"; and exit 1

set CURRENT (env DESK_ENV=$HOME/.desk/desks/hello.fish desk)
echo $CURRENT | grep "say_hello  Args: <hello_to>. Say hello to someone." >/dev/null
test $status -ne 0; and echo "say_hello command not found"; and exit 1

set RAN (desk run hello 'say_hello mrfish')
echo $RAN | grep "Hello mrfish" >/dev/null
test $status -ne 0; and echo "Desk run with 'hello' failed"; and exit 1

set RAN (env DESK_SHELL_ARGS=-l desk run hello 'status --is-login && echo YES || echo NO')
echo $RAN | grep "YES" >/dev/null
test $status -ne 0; and echo "Desk run with DESK_SHELL_ARGS failed"; and exit 1

echo "tests pass."

exit 0
