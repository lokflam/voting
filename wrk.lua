-- env vote_id=1 num=100 wrk -t1 -c1 -s wrk.lua http://18.212.21.72:9009/ballot/cast
code = 1
vote_id = ""

function init(args)
    vote_id = os.getenv("vote_id")
    num = os.getenv("num")
end

request = function()
   local data = [[{
        "vote_id": "%s",
        "choice": "1",
        "code": "%s"
   }]]
   wrk.method = "POST"
   wrk.body   = string.format(data, vote_id, tostring(code))
   wrk.headers["Content-Type"] = "application/json"
   return wrk.format(nil)
end

function response()
   if tostring(code) == num then
      wrk.thread:stop()
   end
   code = code + 1
end