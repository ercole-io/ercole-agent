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
if [ $os_system == 'linux' ] 
	then
	#Retrieve number of threads available
	thread_number=`grep processor /proc/cpuinfo | wc -l`
	#Linux system, sar files at this path /var/log/sa sa01..sa31
	number_of_files=`ls -rt /var/log/sa/sa[0-9][0-9] | wc -l`
	if [ $number_of_files -gt 0 ]
	then		
		for file in `ls -rt /var/log/sa/sa[0-9][0-9]`
		do
			unset sar_file_array
			#Each line contained in the sar file is "saved" in an array (to avoid opening the file n times to read cell by cell)
			while read line
			do
				sar_file_array[${#sar_file_array[@]}]=$line
			done < <(sar -f $file)
			#Loop through the array containing the rows of the sar file just read
			for i in ${!sar_file_array[@]}
			do		
				#For each row you retrieve 4 cells
				stringaprova=`(echo ${sar_file_array[$i]} | awk '{print $1"|"$2"|"$4"|"$9}')`
				IFS="|" read -r -a myarray <<< "$stringaprova"
				#If it is the first row of the array, the date is retrieved.
				if [ $i == 0 ]
				then
					file_date="${myarray[2]}" 
					dates_second=(`date -d "$file_date" +%s`)
				fi
				hour_column="${myarray[0]}" 
				ampm_column="${myarray[1]}" 						
				value_column="${myarray[3]}" 
				#Only rows that have valued the CPU idle field are considered (excluded the row with %idle header)
				if  ! [ -z $hour_column ] && [ $hour_column != "Average:" ] && ! [ -z $value_column ] && [ $value_column != "%idle" ] 
				then
					#Values saved in three arrays: date expressed in seconds / date-time expressed in seconds / cpu idle value
					dates_array+=($dates_second)
					datetime_array+=(`date -d "$file_date $hour_column$ampm_column" +%s`)
					values_array+=($value_column)
				fi				
			done 
		done		
		#Now (time includes) in seconds
		now_saved_time=`date '+%s'`
		#Today (time 00:00:00) in seconds
		now_saved_no_time=`date -d \`date +%Y-%m-%d\` +%s`
	    current_date_no_time_second_changed_month=0
		month_count=0
		day_month_count=0		
		v_cpu_host_month_avg=0
		daily_count_array=(0 0 0 0 0 0 0)
		v_cpu_host_daily_avg_array=(0 0 0 0 0 0 0)
		current_date_no_time_second_changed_week_array=(0 0 0 0)
		day_week_count_array=(0 0 0 0)
		week_count_array=(0 0 0 0)
		v_cpu_host_week_avg_array=(0 0 0 0)	
		v_time_series_array=()
		v_cpu_host_series_avg_array=()
		time_series_count=0
		current_week=-1
		#All arrays have the same number of elements by design (loop by using one)
		for i in ${!values_array[@]}		
		do
			#Difference expressed in hours between now and date of record.
			hours_difference_time=$(( (now_saved_time - ${datetime_array[$i]}) / 3600 )) 
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
				#Counting rows measurements sar in the month
				(( month_count=month_count+1 ))		
				#Total sum of measurements in the month
				v_cpu_host_month_avg=$(echo "scale=2;$v_cpu_host_month_avg+${values_array[$i]}" | bc)
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
				#Count rows measurements sar in current week
				(( week_count_array[$current_week]=week_count_array[$current_week]+1))
				#Total sum of measurements in the current week
				v_cpu_host_week_avg_array[$current_week]=$(echo "scale=2;${v_cpu_host_week_avg_array[$current_week]}+${values_array[$i]}" | bc)	
			fi
							
			#Last 7 days	
			if [ $day_number_no_time -lt 7 ] 
			then
				#Counting measurements sar day by day (array) (index 0 today, index 1 yesterday, ...)
				(( daily_count_array[$day_number_no_time]++ ))
				#Total sum of measurements day by day (array)
				v_cpu_host_daily_avg_array[$day_number_no_time]=$(echo "scale=2;${v_cpu_host_daily_avg_array[$day_number_no_time]}+${values_array[$i]}" | bc)
			fi
			
			#Last day giorno (from now -24 hours)
			#Difference between now and date of the record (/24 -> expressed in days)
			day_number_time=$(( hours_difference_time / 24 ))
			if [ $day_number_time -lt 1 ] 
			then
				#From the date expressed in seconds we get back the date in the format ddmmHH:MM
				v_time_series_array[$time_series_count]=`date -d @"${datetime_array[$i]}" '+%d%m%H:%M'`
				v_cpu_host_series_avg_array[$time_series_count]=${values_array[$i]}
				(( time_series_count=time_series_count+1 ))
			fi			
		done
		
		#Print monthly data (show only if the number of days for which you have sar files is greater/equal than month_soft_limit [arbitrarily decided soft limit])
		if [ $day_month_count -ge $month_soft_limit ]
		then
			echo $(echo "scale=2;(100-($v_cpu_host_month_avg/$month_count))*$thread_number/100" | bc)
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
				echo $(echo "scale=2;(100-(${v_cpu_host_week_avg_array[$i]}/${week_count_array[$i]}))*$thread_number/100" | bc)
			else
				echo "N/A"
			fi
		done
		
		#Show last7 days
		for i in ${!daily_count_array[@]}
		do	
			if [ ${daily_count_array[$i]} -gt 0 ]
			then
				echo $(echo "scale=2;(100-(${v_cpu_host_daily_avg_array[$i]}/${daily_count_array[$i]}))*$thread_number/100" | bc)
			else
				echo "N/A"
			fi
		done
		
		#Show time serie last 24 hours
		if [ ${#v_time_series_array[@]} -gt 0 ]
		then
			for i in ${!v_time_series_array[@]}
			do
				v_time_series_array_string=${v_time_series_array[$i]}\|\|\|$(echo "scale=2;(100-${v_cpu_host_series_avg_array[$i]})*$thread_number/100" | bc)
				echo $v_time_series_array_string
			done
		else
			echo "N/A"
		fi
	#else
		#echo "No files"
	fi
#else
	#echo "No linux"
fi
unset dates_array
unset datetime_array
unset values_array
unset sar_file_array
unset daily_count_array
unset v_cpu_host_daily_avg_array
unset v_cpu_host_series_avg_array
unset v_time_series_array
unset current_date_no_time_second_changed_week_array
unset day_week_count_array
unset week_count_array
unset v_cpu_host_week_avg_array