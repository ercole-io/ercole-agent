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

#Notes

#tps: indicate the number of transfers per second that were issued to the device.
#rd_sec/s: number of sectors read from the device. The size of a sector is 512 bytes.
#wr_sec/s: number of sectors written to the device. The size of a sector is 512 bytes.

#IOMB = (rd_sec/s + wr_sec/s) * 512 / (1024 * 1024 ) -> MB/s

#Needed for AM/PM in the sar output
export LC_ALL=en_US.UTF-8

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
sector_size=512;
#Disks quota limit for calculation (lots of disks -> too much time to check sar files)
maximum_number_of_disks=99;
#Print time series (false->only averages, not last 24 hours time series)
include_time_series=true;
disks_number=0;
if [ $os_system == 'linux' ] 
	then
	#Check number of disks sar retrieved
	today_exists=`ls -l /var/log/sa/sa$today | wc -l`
	if [ $today_exists -gt 0 ]
	then
		disks_number=`sar -dp -f /var/log/sa/sa$today | awk '(/sd/ || /xvd/) && /Average/' | wc -l`
	fi
	#We have at least a sar file and not more than $maximum_number_of_disks
	if [ $today_exists -gt 0 ] && [ $disks_number -le $maximum_number_of_disks ]
	then	
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
				done < <(sar -dp -f $file | awk 'NR==1 || ((/sd/ || /xvd/) && /Average/) || /RESTART/')
				restart_line_number=0
				first_average_for_this_iteration=0
				for i in ${!sar_array_average[@]}
				do
					#Retrieve data from current file
					if [ $i == 0 ]
					then
						file_date=`(echo ${sar_array_average[$i]} | awk '{print $4'})`
						dates_second=(`date -d "$file_date" +%s`)
					else
						#For each file an element in each array will be inserted date, tps (sum over different disks), iomb (rd_sec/s+wr_sec/s sum over different disks)
						stringaprova=`(echo ${sar_array_average[$i]} | awk '{print $3"|"$4"|"$5}')`			
						IFS="|" read -r -a myarray <<< "$stringaprova"
										
						#if it's a restart line continue with the next iteration						
						if [ `(echo ${sar_array_average[$i]} | awk '/RESTART/' | wc -l)` -gt 0 ]
						then
							restart_line_number=$i
							continue
						fi
						
						#If is the first line after date line and not a restart line OR the first line after a restart line
						if [[ $i == 1 && $i -ne $restart_line_number ]] || [[ $i == $(($restart_line_number+1)) ]]
						then
							#Is the first iteration for this file, initialize arrays with the first values
							if [ $first_average_for_this_iteration == 0 ] 
							then
								first_average_for_this_iteration=1
								#Array creation
								dates_array+=($dates_second)
								tps_array+=(${myarray[0]}) 
								iomb_array+=($(echo "scale=2;(${myarray[1]}+${myarray[2]})" | bc))
							#Is an average line after the restart line, arrays have to be reinitialized with the current values
							else
								(( last_index=${#dates_array[@]}-1 ))
								tps_array[$last_index]=$(echo "scale=2;${myarray[0]}" | bc)
								iomb_array[$last_index]=$(echo "scale=2;(${myarray[1]}+${myarray[2]})" | bc)
							fi				
						else					
							#Other lines to be summarized to the initialized line
							(( last_index=${#dates_array[@]}-1 ))
							tps_array[$last_index]=$(echo "scale=2;${tps_array[$last_index]}+${myarray[0]}" | bc)
							iomb_array[$last_index]=$(echo "scale=2;${iomb_array[$last_index]}+((${myarray[1]}+${myarray[2]}))" | bc)
						fi														
					fi
				done
			done
			
			if [ $include_time_series = true ]
			then
				for file in `ls -rt /var/log/sa/sa$today /var/log/sa/sa$yesterday`
				do
					unset sar_array_average
					#Each line contained in the sar files (considering the awk filter) is "saved" in an array (to avoid opening the file n times to read cell by cell)
					while read line
					do
						sar_array_average[${#sar_array_average[@]}]=$line
					#Retrieve only first row and average rows for only sd* or xvd* disks
					done < <(sar -dp -f $file | awk 'NR==1 || ((/sd/ || /xvd/) && !/Average/)')
					for i in ${!sar_array_average[@]}
					do
						#Retrieve data from current file
						if [ $i == 0 ]
						then
							file_date=`(echo ${sar_array_average[$i]} | awk '{print $4'})`
						else
							#For each line with different time tps and iomb will we saved (summarized for different disks)
							stringaprova=`(echo ${sar_array_average[$i]} | awk '{print $1"|"$2"|"$4"|"$5"|"$6}')`			
							IFS="|" read -r -a myarray <<< "$stringaprova"
							hour_column="${myarray[0]}" 
							ampm_column="${myarray[1]}" 
							tps_column="${myarray[2]}" 
							iomb_column=$(echo "scale=2;(${myarray[3]}+${myarray[4]})" | bc) 		
							#Current line datetime in seconds
							current_line_datetime=(`date -d "$file_date $hour_column$ampm_column" +%s`)
							#if last element datetime = current element datetime -> sum tps and iomb (two or more different disks -> sum values)
							#only "one record" per different datetime
							if [ ${#datetime_array_daily[@]} -gt 0 ] 
							then 					
								if [ ${datetime_array_daily[-1]} == $current_line_datetime ]
								then 
									#Last inserted element index
									(( last_index=${#datetime_array_daily[@]}-1 ))
									#Datetime already exist -> sum tps e iomb with previous record 
									tps_array_daily[$last_index]=$(echo "scale=2;${tps_array_daily[$last_index]}+$tps_column" | bc)
									iomb_array_daily[$last_index]=$(echo "scale=2;${iomb_array_daily[$last_index]}+$iomb_column" | bc)
								else
									#Datetime not exist -> new record 
									#Values saved in three arrays: date-time expressed in seconds / tps value / iomb value
									datetime_array_daily+=($current_line_datetime)
									tps_array_daily+=($tps_column)
									iomb_array_daily+=($(echo "scale=2;$iomb_column" | bc))
								fi		
							else
								#Empty array -> new record 
								#Values saved in three arrays: date-time expressed in seconds / tps value / iomb value
								datetime_array_daily+=($current_line_datetime)
								tps_array_daily+=($tps_column)
								iomb_array_daily+=($(echo "scale=2;iomb_column" | bc))
							fi		
									
						fi
					done
				done
			fi
					
			
			current_date_no_time_second_changed_month=0
			month_count=0
			day_month_count=0		
			v_tps_month_avg=0
			v_iomb_month_avg=0
			daily_count_array=(0 0 0 0 0 0 0)
			v_tps_daily_avg_array=(0 0 0 0 0 0 0)
			v_iomb_daily_avg_array=(0 0 0 0 0 0 0)
			current_date_no_time_second_changed_week_array=(0 0 0 0)
			day_week_count_array=(0 0 0 0)
			week_count_array=(0 0 0 0)
			v_tps_week_avg_array=(0 0 0 0)	
			v_iomb_week_avg_array=(0 0 0 0)	
			current_week=-1
			v_time_series_array=()
			v_tps_series_avg_array=()
			v_iomb_series_avg_array=()
			time_series_count=0
			
			#Montly/weekly/daily averages evaluation
			#All arrays have the same number of elements by design (loop by using one)
			for i in ${!tps_array[@]}		
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
					#Counting rows measurements sar in the month
					(( month_count=month_count+1 ))	
					#Total sum of measurements in the month
					v_tps_month_avg=$(echo "scale=2;$v_tps_month_avg+${tps_array[$i]}" | bc)
					v_iomb_month_avg=$(echo "scale=2;$v_iomb_month_avg+${iomb_array[$i]}" | bc)
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
					v_tps_week_avg_array[$current_week]=$(echo "scale=2;${v_tps_week_avg_array[$current_week]}+${tps_array[$i]}" | bc)	
					v_iomb_week_avg_array[$current_week]=$(echo "scale=2;${v_iomb_week_avg_array[$current_week]}+${iomb_array[$i]}" | bc)	
				fi
								
				#Last 7 days	
				if [ $day_number_no_time -lt 7 ] 
				then
					#Counting measurements sar day by day (array) (index 0 today, index 1 yesterday, ...)
					(( daily_count_array[$day_number_no_time]++ ))
					#Total sum of measurements day by day (array)
					v_tps_daily_avg_array[$day_number_no_time]=$(echo "scale=2;${v_tps_daily_avg_array[$day_number_no_time]}+${tps_array[$i]}" | bc)
					v_iomb_daily_avg_array[$day_number_no_time]=$(echo "scale=2;${v_iomb_daily_avg_array[$day_number_no_time]}+${iomb_array[$i]}" | bc)
				fi		
			done
			
			#Last 24 hours averages
			for i in ${!tps_array_daily[@]}		
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
					v_tps_series_avg_array[$time_series_count]=${tps_array_daily[$i]}
					v_iomb_series_avg_array[$time_series_count]=${iomb_array_daily[$i]}
					(( time_series_count=time_series_count+1 ))
				fi		
			done
					
			#Print monthly data (show only if the number of days for which you have sar files is greater/equal than month_soft_limit [arbitrarily decided soft limit])
			if [ $day_month_count -ge $month_soft_limit ]
			then
				echo $(echo "scale=2;$v_tps_month_avg/$month_count" | bc)\|\|\|$(echo "scale=2;($v_iomb_month_avg/$month_count)*$sector_size/(1024*1024)" | bc)
			else
				echo "N/A|||N/A"
			fi
			
			#Output last 4 weeks averages and maximums
			#Loops over weeks
			for i in ${!day_week_count_array[@]}
			do
				#Weekly output (show only if the number of days for which you have sar files is greater/equal than week_soft_limit [arbitrarily decided soft limit])
				if [ ${day_week_count_array[$i]} -ge $week_soft_limit ]
				then
					echo $(echo "scale=2;${v_tps_week_avg_array[$i]}/${week_count_array[$i]}" | bc)\|\|\|$(echo "scale=2;(${v_iomb_week_avg_array[$i]}/${week_count_array[$i]})*$sector_size/(1024*1024)" | bc)
				else
					echo "N/A|||N/A"
				fi
			done
			
			#Show last7 days
			for i in ${!daily_count_array[@]}
			do	
				if [ ${daily_count_array[$i]} -gt 0 ]
				then
					echo $(echo "scale=2;${v_tps_daily_avg_array[$i]}/${daily_count_array[$i]}" | bc)\|\|\|$(echo "scale=2;(${v_iomb_daily_avg_array[$i]}/${daily_count_array[$i]})*$sector_size/(1024*1024)" | bc)
				else
					echo "N/A|||N/A"
				fi
			done		
	
			#Show time serie last 24 hours
			if [ ${#v_time_series_array[@]} -gt 0 ]
			then
				for i in ${!v_time_series_array[@]}
				do
					v_time_series_array_string=${v_time_series_array[$i]}\|\|\|$(echo "scale=2;${v_tps_series_avg_array[$i]}" | bc)\|\|\|$(echo "scale=2;(${v_iomb_series_avg_array[$i]})*$sector_size/(1024*1024)" | bc)
					echo $v_time_series_array_string
				done
			else
				echo "N/A|||N/A"
			fi
		fi
	fi
fi
unset sar_array_average
unset dates_array
unset tps_array
unset iomb_array
unset datetime_array_daily
unset tps_array_daily
unset iomb_array_daily
unset daily_count_array
unset v_tps_daily_avg_array
unset v_iomb_daily_avg_array
unset current_date_no_time_second_changed_week_array
unset day_week_count_array
unset week_count_array
unset v_tps_week_avg_array
unset v_iomb_week_avg_array
unset v_time_series_array
unset v_tps_series_avg_array
unset v_iomb_series_avg_array