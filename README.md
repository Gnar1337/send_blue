# send_blue

<h2>Steps to run:</h2>

1. Clone the Repo `git clone https://github.com/Gnar1337/send_blue.git`
2. `cd send_blue`
3. change `QUEUE_INTERVAL` value to change the interval between sending messages in the `docker-compose.yaml`
4. run `docker compose up` 
5. Navigate to `http://localhost:5173/` to schedule
6. to change the queue send time use this url get `http://localhost:8080/gateway/interval?seconds=180` this example will change the queue time to 180 seconds

<h2>stack</h2>
<ul>
<li>vue</li>
<li>go</li>
<li>postgres</li>
<li>gRPC</li>
<li>Docker</li>
</ul>

Gnar1337
