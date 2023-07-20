#!/bin/bash

# Copyright (c) 2023 Sorint.lab S.p.A.
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

#Retrieve os (linux, aix, sunos)
inizio=`date`
os_system=`uname -a | awk '{print $1}' | tr '[:upper:]' '[:lower:]'`
month_soft_limit=15;
week_soft_limit=4;
today=`date '+%d'`
yesterday=`date --date="1 day ago" '+%d'`
#Now (time includes) in seconds
now_saved_time=`date '+%s'`
#Today (time 00:00:00) in seconds
now_saved_no_time=`date -d \`date +%Y-%m-%d\` +%s`
if [ $os_system == 'linux' ] 
	then
	#Retrieve number of threads available
	thread_number=`grep processor /proc/cpuinfo | wc -l`
	#Linux system, sar files at this path /var/log/sa sa01..sa31
	number_of_files=`ls -rt /var/log/sa/sa[0-9][0-9] | wc -l`
	if [ $number_of_files -gt 0 ]
	then
		#Scan all files for montly, weekhly e daily averages
		for file in `ls -rt /var/log/sa/sa[0-9][0-9]`
		do
			unset sar_array_average
			#Each line contained in the sar files (considering the awk filter) is "saved" in an array (to avoid opening the file n times to read cell by cell)
			while read line
			do
				sar_array_average[${#sar_array_average[@]}]=$line
			#Retrieve only first row and average rows for only sd* or xvd* disks
			done < <(sar -f $file | awk 'NR==1 || /Average/')
			for i in ${!sar_array_average[@]}
			do
				#Retrieve data from current file
				if [ $i == 0 ]
				then
					file_date=`(echo ${sar_array_average[$i]} | awk '{print $4'})`
					dates_second=(`date -d "$file_date" +%s`)
				else
					#For each file an element in each array will be inserted date, tps (sum over different disks), iomb (rd_sec/s+wr_sec/s sum over different disks)
					cpu_column=`(echo ${sar_array_average[$i]} | awk '{print $8'})`	
					#Line 1, arrays created
					dates_array+=($dates_second)
					cpu_array+=($cpu_column) 												
				fi
			done
		done
		
		for file in `ls -rt /var/log/sa/sa$today /var/log/sa/sa$yesterday`
		do
			unset sar_array_average
			#Each line contained in the sar files (considering the awk filter) is "saved" in an array (to avoid opening the file n times to read cell by cell)
			while read line
			do
				sar_array_average[${#sar_array_average[@]}]=$line
			#Retrieve only first row and rows for only sd* or xvd* disks
			done < <(sar -f $file | awk 'NR==1 || (!/Average/ && !/RESTART/ && !/CPU/)')
			for i in ${!sar_array_average[@]}
			do
				#Retrieve data from current file
				if [ $i == 0 ]
				then
					file_date=`(echo ${sar_array_average[$i]} | awk '{print $4'})`
				else
					#For each line with different time tps and iomb will we saved (summarized for different disks)
					stringaprova=`(echo ${sar_array_average[$i]} | awk '{print $1"|"$2"|"$9}')`			
					IFS="|" read -r -a myarray <<< "$stringaprova"
					hour_column="${myarray[0]}" 
					ampm_column="${myarray[1]}" 
					cpu_column="${myarray[2]}"  		
					#Current line datetime in seconds
					current_line_datetime=(`date -d "$file_date $hour_column$ampm_column" +%s`)
					#Only rows that have valued the CPU idle field are considered (excluded the row with %idle header)
					if  ! [ -z $hour_column ] && ! [ -z $ampm_column ] && ! [ -z $cpu_column ] 
					then
						#Values saved in three arrays: date expressed in seconds / date-time expressed in seconds / cpu idle value
						datetime_array_daily+=(`date -d "$file_date $hour_column$ampm_column" +%s`)
						cpu_array_daily+=($cpu_column)
					fi										
				fi
			done
		done
			
		current_date_no_time_second_changed_month=0
		day_month_count=0		
		v_cpu_month_avg=0
		daily_count_array=(0 0 0 0 0 0 0)
		v_cpu_daily_avg_array=(0 0 0 0 0 0 0)
		current_date_no_time_second_changed_week_array=(0 0 0 0)
		day_week_count_array=(0 0 0 0)
		v_cpu_week_avg_array=(0 0 0 0)	
		current_week=-1
		v_time_series_array=()
		v_cpu_series_avg_array=()
		time_series_count=0
		
		#Montly/weekly/daily averages evaluation
		#All arrays have the same number of elements by design (loop by using one)
		for i in ${!cpu_array[@]}		
		do
			#Difference expressed in hours between now and date of the record (considering only the day and not the time of the records)
			hours_difference_no_time=$(( (now_saved_no_time - ${dates_array[$i]}) / 3600 )) 
			#Difference between now and date (no time) of the record (/24 -> expressed in days)			
			day_number_no_time=$(( hours_difference_no_time / 24 ))
			
			#Last 30 days (montly average)
			if [ $day_number_no_time -le 30 ] 
			then			
				#Counter to know how many effective days have sar files (it could be the case that in 30 days interval, sar files are not available for x days?)
				current_date_no_time_seconds_month=${dates_array[$i]}
				if [ $current_date_no_time_second_changed_month != $current_date_no_time_seconds_month ] 
				then 
					current_date_no_time_second_changed_month=$current_date_no_time_seconds_month
					(( day_month_count=day_month_count+1 ))
				fi			
				#Total sum of measurements in the month
				v_cpu_month_avg=$(echo "scale=2;$v_cpu_month_avg+${cpu_array[$i]}" | bc)
			fi
			
			#Integer that represents the week number (0 current, 1 next week, ..., 3)
			current_week=$(($day_number_no_time / 7))
			if [ $current_week -lt 4 ] 
			then				
				#Counter to know how many effective days have sar files (it could be the case that in 7 days interval, sar files are not available for x days?)
				current_date_no_time_seconds_week=${dates_array[$i]}
				if [ ${current_date_no_time_second_changed_week_array[$current_week]} != $current_date_no_time_seconds_week ] 
				then 
					current_date_no_time_second_changed_week_array[$current_week]=$current_date_no_time_seconds_week
					(( day_week_count_array[$current_week]=day_week_count_array[$current_week]+1))
				fi	
				#Total sum of measurements in the current week
				v_cpu_week_avg_array[$current_week]=$(echo "scale=2;${v_cpu_week_avg_array[$current_week]}+${cpu_array[$i]}" | bc)		
			fi
							
			#Last 7 days	
			if [ $day_number_no_time -lt 7 ] 
			then
				#Counting measurements sar day by day (array) (index 0 today, index 1 yesterday, ...)(should be 1 [only averages retrieved] or 0, no sar file for that day)
				(( daily_count_array[$day_number_no_time]++ ))
				#Total sum of measurements day by day (array)
				v_cpu_daily_avg_array[$day_number_no_time]=${cpu_array[$i]}
			fi		
		done
		
		#Last 24 hours averages
		for i in ${!cpu_array_daily[@]}		
		do
			#Difference expressed in hours between now and date of record.
			hours_difference_time=$(( (now_saved_time - ${datetime_array_daily[$i]}) / 3600 )) 
			#Difference between now and date of the record (/24 -> expressed in days)
			day_number_time=$(( hours_difference_time / 24 ))
			#Last day (from now -24 hours)		
			if [ $day_number_time -lt 1 ] 
			then
				#From the date expressed in seconds we get back the date in the format ddmmHH:MM
				v_time_series_array[$time_series_count]=`date -d @"${datetime_array_daily[$i]}" '+%d%m%H:%M'`
				v_cpu_series_avg_array[$time_series_count]=${cpu_array_daily[$i]}
				(( time_series_count=time_series_count+1 ))
			fi		
		done
		
		#Print monthly data (show only if the number of days for which you have sar files is greater/equal than month_soft_limit [arbitrarily decided soft limit])
		if [ $day_month_count -ge $month_soft_limit ]
		then
			echo $(echo "scale=2;(100-($v_cpu_month_avg/$day_month_count))*$thread_number/100" | bc)
		else
			echo "N/A"
		fi
		
		#Output last 4 weeks averages and maximums
		#Loops over weeks
		for i in ${!day_week_count_array[@]}
		do
			#Weekly output (show only if the number of days for which you have sar files is greater/equal than week_soft_limit [arbitrarily decided soft limit])
			if [ ${day_week_count_array[$i]} -ge $week_soft_limit ]
			then
				echo $(echo "scale=2;(100-(${v_cpu_week_avg_array[$i]}/${day_week_count_array[$i]}))*$thread_number/100" | bc)
			else
				echo "N/A"
			fi
		done
		
		#Show last7 days
		for i in ${!daily_count_array[@]}
		do	
			if [ ${daily_count_array[$i]} -gt 0 ]
			then
				echo $(echo "scale=2;(100-${v_cpu_daily_avg_array[$i]})*$thread_number/100" | bc)
			else
				echo "N/A"
			fi
		done		

		#Show time serie last 24 hours
		if [ ${#v_time_series_array[@]} -gt 0 ]
		then
			for i in ${!v_time_series_array[@]}
			do
				v_time_series_array_string=${v_time_series_array[$i]}\|\|\|$(echo "scale=2;(100-${v_cpu_series_avg_array[$i]})*$thread_number/100" | bc)
				echo $v_time_series_array_string
			done
		else
			echo "N/A"
		fi
	fi
fi
unset sar_array_average
unset dates_array
unset cpu_array
unset datetime_array_daily
unset cpu_array_daily
unset daily_count_array
unset v_cpu_daily_avg_array
unset current_date_no_time_second_changed_week_array
unset day_week_count_array
unset week_count_array
unset v_cpu_week_avg_array
unset v_time_series_array
unset v_cpu_series_avg_array