workspace "gpc"

configurations {"Debug", "Release"}

project "gpc"
    kind "ConsoleApp"
    language "C"

    targetdir "bin/%{cfg.buildcfg}"

    files { "**.h", "**.c" }
 
    filter "configurations:Debug"
       defines { "DEBUG" }
       symbols "On"
 
    filter "configurations:Release"
       defines { "NDEBUG" }
       optimize "On"