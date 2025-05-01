Set FontSize 20
Set Width 1200
Set Height 600
Set BorderRadius 40
Set WindowBar Colorful
Set WindowBarSize 50
Output wish-auth.gif

Hide
Type "tmux kill-session -t demo" Enter
Type "tmux new -A -s demo" Enter
Type "tmux set status off" Enter
Type "bash --rcfile rc.bash" Enter
Type "tmux splitw bash --rcfile rc.bash" Enter
Ctrl+B
Type k
Type "clear" Enter
Show

Type "go run ./wish-auth/main.go"
Enter
Sleep 1s

Ctrl+B
Type j

Sleep 1s
Type "ssh -p 23234 -o PreferredAuthentications=publickey localhost" Enter
Sleep 5s
Enter

Type "ssh -p 23234 -o PreferredAuthentications=password localhost" Enter
Sleep 2s
Type "how you turn this on" Enter
Sleep 5s
Enter

Type "ssh -p 23234 -o PreferredAuthentications=keyboard-interactive localhost" Enter
Sleep 2s
Type "5" Enter
Sleep 2s
Type "vim" Enter
Sleep 2s
Type "nope" Enter
Sleep 10s
Enter

Hide
Type "tmux kill-session -t demo" Enter
Show
