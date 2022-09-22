Write-Output 'Welcome to the NWRS Client, First input your server IP.'
$Server = Read-Host -Prompt 'Input your server IP'
For (;;) {
    Write-Output 'Choose the number of the action you would like to perform.'
Write-Output "Create a user? [1]`nRemove a user? [2]`ncreate a webserver for a existing user? [3]`nremove an existing website belonging to a user? [4]`nExit [5]"
$Choice = Read-Host -Prompt 'Option'
switch ($Choice) {
        1 {
            $User = Read-Host -Prompt 'Input the user name'
            $Passw = Read-Host -Prompt 'Input the user password'
            Invoke-RestMethod -Uri ("http://" + $Server + ":1234/wrs/user?user=" + $User + "&pass=" + $Passw) -Method POST
        }
        2 {
            $User = Read-Host -Prompt 'Input the user name'
            Invoke-RestMethod -Uri ("http://" + $Server + ":1234/wrs/user?user=" + $User) -Method DELETE
        }
        3 {
            $User = Read-Host -Prompt 'Input the user name'
            Invoke-RestMethod -Uri ("http://" + $Server + ":1234/wrs/container?user=" + $User) -Method POST
        }
        4 {
            $User = Read-Host -Prompt 'Input the user name'
            Invoke-RestMethod -Uri ("http://" + $Server + ":1234/wrs/container?user=" + $User) -Method DELETE
        }
        5 {
            Exit
        }
    }
}
