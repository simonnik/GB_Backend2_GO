--[[

    4 cases:

    1) incorrect acl token
    2) successful acl check and add data to DB
    3) succsessful / fail DB query row
    4) successful list all data from DB

        http://tylerneylon.com/a/learn-lua/

--]]

math.randomseed(os.time())
request = function()
    local k = math.random(0, 1000)
    local r = math.random(0, 1000)
    local t
    local url
    local method

    -- expect error on acl-check
    if k > 950 then
        t = "incorrect_admin_token"
        url = "/new-entity?token="..t.."&id="..k.."&data="..k
        method = "POST"

    -- expect successful acl-check and add data to DB
    elseif k > 300 then
        if r > 500 then
            t = "admin_secret_token"
            url = "/new-entity?token="..t.."&id="..k.."&data="..k
            method = "POST"
        else
            url = "/get-entity?id="..k
            method = "GET"
        end

    -- expect successful list all data from DB
    else
        t = "admin_secret_token"
        url = "/entities"
        method = "GET"
    end

    return wrk.format(method, url)
end