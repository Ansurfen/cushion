---@diagnostic disable: undefined-global
Import({ "cushion" })

local tui = require("cushion-tui")

function TestFancyList()
    local choice = tui.FancyList({
        { title = "cushion", detail = "cushion is base for hulo" }
    })
    print(choice.ChoiceTitle, choice.ChoiceDetail)
end

function TestSpinner()
    tui.Spinner(function()
        local sleep = function(n)
            local t = os.time()
            repeat
            until os.time() > t + n
        end

        local setInterval = function(callback, interval)
            sleep(interval)
            callback()
        end
        print("start")
        setInterval(function()
            print("Hello World!")
        end, 5)
        print("finish")
    end)
end

function TestSimpleList()
    local choice = tui.SimpleList({
        title = "cushion",
        choices = { "a", "b", "c" }
    })
    print(choice)
end

function TestTextInput()
    local res = tui.TextInput({
        { name = "host" },
        { name = "password", echomode = true }
    })
    print(res[1], res[2])
end

function TestMultiSelect()
    local res = tui.MultiSelect({
        title = "cushion",
        choices = { "a", "b", "c" }
    })
    print(res[1])
end

function TestBatchSpinner()
    tui.BatchSpinner({
        {
            name = "fetch",
            callback = function()

            end
        },
        {
            name = "unzip",
            callback = function ()
                
            end
        }
    })
end
