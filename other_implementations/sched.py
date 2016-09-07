# Jack Schefer, Began 6/15/16
#
# This implementation scrapes the current day's schedule from the website instead of the api.
#
## Purpose: get the schedule from the TJ Intranet in order and print 
#               it out to the terminal
#
from subprocess import call
from datetime   import datetime
#
def main():
   #
   print()
   #
   # 1. copy the html from ion into a file called sched.txt
   args = 'curl -s https://ion.tjhsst.edu > sched.txt'
   call(args, shell = True)
   #
   # 2. read the file into an array of lines
   fname = 'sched.txt'
   lines = []
   with open(fname, 'r') as afile:
      lines = afile.readlines()
   #
   # 3. delete teh file sched.txt
   args = 'rm ' + fname
   call(args, shell = True)
   #
   # 4. find and print the date, day type, and current time
   for l in lines:
      #
      if 'class=\"schedule-date\"' in l: print(innerHTML(l))
      if 'class=\"day-name'   in l: print(innerHTML(l))
      #
   #
   strtime = datetime.now().time().isoformat()
   hr = int(strtime[0 : 2])
   mn = int(strtime[3 : 5])
   suffix = 'am'
   #
   if hr > 12: 
      hr -= 12
      suffix = 'pm'
   #
   print('Time:', str(hr) + ':' + str(mn), suffix)
   #
   print()
   #
   # 5. search for the blocks within the <th> tags
   removing = []
   for l in lines:
      #
      strt = l[:3]
      if strt != '<th' and strt != '<td':removing.append(l) 
      #
   #
   for r in removing:
      #
      lines.remove(r)
      #
   #
   # 6. print out each pair of header and times, make the current block bold
   for i in range(len(lines) // 2):
      #
      header = lines[2 * i]
      str1 = innerHTML(header)
      #
      times = lines[2 * i + 1]
      str2 = innerHTML(times)
      #
      if inInterval(str2):
        print(Color.CYAN + Color.BOLD +  str1, '  \t', str2 +  Color.END + Color.END)
      #
      else:
        print(str1 + '  ', '\t', str2)
   #
   print()
   #
#
#################################################################################
#
def inInterval(time_string):    # time string is of the format hh:mm - hh:mm (or 1 h)
    #
    # 1. append a zero to the times if necessary
    firstColon  = time_string.find(':')
    secondColon = time_string.rfind(':')
    #
    if int(time_string[firstColon - 1]) >= 6: 
        time_string = time_string[: firstColon - 1] + '0' + time_string[firstColon - 1 :]
        secondColon += 1
        firstColon  += 1
    #
    hour1 = int(time_string[firstColon - 2  : firstColon    ])
    min1  = int(time_string[firstColon + 1  : firstColon + 3])
    #
    if int(time_string[secondColon - 1]) >= 6: 
        time_string = time_string[: secondColon - 1] + '0' + time_string[secondColon - 1]
        secondColon += 1
    #
    hour2 = int(time_string[secondColon - 2 : secondColon    ])
    min2  = int(time_string[secondColon + 1 : secondColon + 3])
    #
    # 2. get the current hour and minute from python
    string_form = datetime.now().time().isoformat()
    hour = int(string_form[0 : 2])
    minn = int(string_form[3 : 5])
    #
    #print('begin: ', hour1, ':', min1)
    #print('curr: ', hour, ':', minn)
    #print('end: ', hour2, ':', min2)
    #
    # 3. compare the two
    return 60 * hour1 + min1 <= 60 * hour + minn <= 60 * hour2 + min2
    #
#
#################################################################################
#
def innerHTML(s):
   #
   pos1 = 0
   for j in range(len(s)):
       if s[j] == '>': 
           pos1 = j
           break
   #
   pos2 = 0
   for j in range(pos1, len(s)):
       if s[j] == '<':
           pos2 = j
           break
   #
   str1 = s[pos1 + 1 : pos2]
   return str1
   #
#
##################################################################################
#
class Color:
   #
   PURPLE    = '\033[95m'
   CYAN      = '\033[96m'
   DARKCYAN  = '\033[36m'
   BLUE      = '\033[94m'
   GREEN     = '\033[92m'
   YELLOW    = '\033[93m'
   RED       = '\033[91m'
   BOLD      = '\033[1m'
   UNDERLINE = '\033[4m'
   END       = '\033[0m'
   #
#
#################################################################################
#
if  __name__ == '__main__':
    #
    main()
    #
#
# End of File.
