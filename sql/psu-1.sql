-- Copyright (c) 2019 Sorint.lab S.p.A.

-- This program is free software: you can redistribute it and/or modify
-- it under the terms of the GNU General Public License as published by
-- the Free Software Foundation, either version 3 of the License, or
-- (at your option) any later version.

-- This program is distributed in the hope that it will be useful,
-- but WITHOUT ANY WARRANTY; without even the implied warranty of
-- MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
-- GNU General Public License for more details.

-- You should have received a copy of the GNU General Public License
-- along with this program.  If not, see <http://www.gnu.org/licenses/>.

VARIABLE STATUS varchar2(100);
VARIABLE PSU_DATE varchar2(100);
VARIABLE DESCRIPTION varchar2(100);
VARIABLE VERSION varchar2(100);
VARIABLE EXTENDVERSION varchar2(100);
VARIABLE EXIST number;
VARIABLE BP varchar2(100);

set lines 8000 pages 0 feedback off verify off
set colsep "|||"
alter session set NLS_DATE_FORMAT='YYYY-MM-DD';


BEGIN
SELECT DBMS_DB_VERSION.VERSION || '.' || DBMS_DB_VERSION.RELEASE into :VERSION FROM v$instance;
SELECT version into :EXTENDVERSION from v$instance;
select count(*) into :EXIST  from  registry$history;
select * into :BP from (select COMMENTS from  registry$history order by action_time DESC) where rownum=1;
-- 11.2
 IF ( :EXTENDVERSION = '11.2.0.4.0' AND :EXIST > 0 AND :BP not like '%BP%') THEN 
        with PSU as
                (
                        select max(COMMENTS) as COMMENTS
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                ),
                DATA as
                (
                        select 
                        case
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 1 THEN TO_DATE('140115','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 2 THEN TO_DATE('140415','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 3 THEN TO_DATE('140714','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 4 THEN TO_DATE('141014','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 5 THEN TO_DATE('150120','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 6 THEN TO_DATE('150414','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 7 THEN TO_DATE('150714','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 8 THEN TO_DATE('151020','YYMMDD')
                                ELSE
                                         TO_DATE(substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1), 'YYMMDD')
                        END as PSU_DATE
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                ),
                STATE as
                (
                        select 
                        case   
                                WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
                                ELSE 'OK' 
                        END as "STATUS"
                        from DATA
                )
                select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS 
                into :DESCRIPTION,:PSU_DATE,:STATUS 
                from PSU,DATA,STATE;
         ELSIF ( :EXTENDVERSION = '11.2.0.4.0' AND :EXIST > 0 AND :BP like '%BP%') THEN
         with PSU as
         (
                        select * from (select COMMENTS
                        from  registry$history
                        where ACTION='APPLY'
                        and COMMENTS like 'BP%'
                        order by action_time DESC) where rownum=1
        ),
         DATA as
         (
                select * from (
            select
             case
                 WHEN COMMENTS = 'BP20' THEN TO_DATE('151015','YYMMDD')
				 WHEN COMMENTS = 'BP19' THEN TO_DATE('150915','YYMMDD')
				 WHEN COMMENTS = 'BP18' THEN TO_DATE('150815','YYMMDD')
				 WHEN COMMENTS = 'BP17' THEN TO_DATE('150715','YYMMDD')
				 WHEN COMMENTS = 'BP16' THEN TO_DATE('150415','YYMMDD')
				 WHEN COMMENTS = 'BP15' THEN TO_DATE('150115','YYMMDD')
				 WHEN COMMENTS = 'BP14' THEN TO_DATE('151214','YYMMDD')
				 WHEN COMMENTS = 'BP13' THEN TO_DATE('151114','YYMMDD')
				 WHEN COMMENTS = 'BP12' THEN TO_DATE('151014','YYMMDD')
				 WHEN COMMENTS = 'BP11' THEN TO_DATE('150914','YYMMDD')
				 WHEN COMMENTS = 'BP10' THEN TO_DATE('150814','YYMMDD')
				 WHEN COMMENTS = 'BP9'  THEN TO_DATE('150714','YYMMDD')
				 WHEN COMMENTS = 'BP8'  THEN TO_DATE('150614','YYMMDD')
				 WHEN COMMENTS = 'BP7'  THEN TO_DATE('150514','YYMMDD')
				 WHEN COMMENTS = 'BP6'  THEN TO_DATE('150414','YYMMDD')
				 WHEN COMMENTS = 'BP5'  THEN TO_DATE('150314','YYMMDD')
				 WHEN COMMENTS = 'BP4'  THEN TO_DATE('150214','YYMMDD')
				 WHEN COMMENTS = 'BP3'  THEN TO_DATE('150114','YYMMDD')
				 WHEN COMMENTS = 'BP2'  THEN TO_DATE('151213','YYMMDD')
				 WHEN COMMENTS = 'BP1'  THEN TO_DATE('151113','YYMMDD')
                 ELSE
                    TO_DATE(substr(COMMENTS,3,7), 'YYMMDD')
             END as PSU_DATE
             from  registry$history
             where ACTION='APPLY'
             and COMMENTS like 'BP%' order by action_time DESC) where rownum=1
         ),
         STATE as
         (
             select
             case
                 WHEN PSU_DATE < SYSDATE-180 THEN 'KO'
                 ELSE 'OK'
             END as "STATUS"
             from DATA
         )
         select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS
         into :DESCRIPTION,:PSU_DATE,:STATUS
         from PSU,DATA,STATE;
   ELSIF ( :EXTENDVERSION = '11.2.0.3.0' AND :EXIST > 0 AND :BP not like '%BP%') THEN 
        with PSU as
                (
                        select max(COMMENTS) as COMMENTS
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                        or COMMENTS like '%DATABASE PATCH SET UPDATE%'
                ),
                DATA as
                (
                        select 
                        case
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 15 THEN TO_DATE('150714','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 14 THEN TO_DATE('150414','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 13 THEN TO_DATE('150120','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 12 THEN TO_DATE('141014','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 11 THEN TO_DATE('140717','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 10 THEN TO_DATE('140415','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 9  THEN TO_DATE('140114','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 8  THEN TO_DATE('131015','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 7  THEN TO_DATE('130716','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 6  THEN TO_DATE('130417','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 5  THEN TO_DATE('130114','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 4  THEN TO_DATE('121015','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 3  THEN TO_DATE('120716','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 2  THEN TO_DATE('120416','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 1  THEN TO_DATE('120116','YYMMDD')
                        END as PSU_DATE
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                        or COMMENTS like '%DATABASE PATCH SET UPDATE%'
                ),
                STATE as
                (
                        select 
                        case   
                                WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
                                ELSE 'OK' 
                        END as "STATUS"
                        from DATA
                )
                select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS 
                into :DESCRIPTION,:PSU_DATE,:STATUS 
                from PSU,DATA,STATE; 
         ELSIF ( :EXTENDVERSION = '11.2.0.3.0' AND :EXIST > 0 AND :BP like '%BP%') THEN
         with PSU as
         (
             select max(COMMENTS) as COMMENTS
             from  registry$history
             where ACTION='APPLY'
             and COMMENTS like 'BP%'
        ),
         DATA as
         (
            select
             case
                 WHEN MAX(COMMENTS) = 'BP28' THEN TO_DATE('150715','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP27' THEN TO_DATE('150415','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP26' THEN TO_DATE('150115','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP25' THEN TO_DATE('151014','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP24' THEN TO_DATE('150714','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP23' THEN TO_DATE('150414','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP22' THEN TO_DATE('150114','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP21' THEN TO_DATE('151013','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP20' THEN TO_DATE('150713','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP19' THEN TO_DATE('150613','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP18' THEN TO_DATE('150513','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP17' THEN TO_DATE('150413','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP16' THEN TO_DATE('150313','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP15' THEN TO_DATE('150213','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP14' THEN TO_DATE('150113','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP13' THEN TO_DATE('151212','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP12' THEN TO_DATE('151112','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP11' THEN TO_DATE('151012','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP10' THEN TO_DATE('150912','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP9'  THEN TO_DATE('150812','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP8'  THEN TO_DATE('150712','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP7'  THEN TO_DATE('150612','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP6'  THEN TO_DATE('150512','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP5'  THEN TO_DATE('150412','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP4'  THEN TO_DATE('150312','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP3'  THEN TO_DATE('150212','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP2'  THEN TO_DATE('150112','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP1'  THEN TO_DATE('151211','YYMMDD')
                 ELSE
                    TO_DATE(substr(COMMENTS,3,7), 'YYMMDD')			 
             END as PSU_DATE
             from  registry$history
             where ACTION='APPLY'
             and COMMENTS like 'BP%'
         ),
         STATE as
         (
             select
             case
                 WHEN PSU_DATE < SYSDATE-180 THEN 'KO'
                 ELSE 'OK'
             END as "STATUS"
             from DATA
         )
         select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS
         into :DESCRIPTION,:PSU_DATE,:STATUS
         from PSU,DATA,STATE;
   ELSIF ( :EXTENDVERSION = '11.2.0.2.0'  AND :EXIST > 0  AND :BP not like '%BP%') THEN 
        with PSU as
                (
                        select max(COMMENTS) as COMMENTS
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                        or COMMENTS like '%DATABASE PATCH SET UPDATE%'
                ),
                DATA as
                (
                        select 
                        case
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 1 THEN TO_DATE('110118','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 2 THEN TO_DATE('110512','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 3 THEN TO_DATE('110811','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 4 THEN TO_DATE('111017','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 5 THEN TO_DATE('120120','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 6 THEN TO_DATE('120417','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 7 THEN TO_DATE('120719','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 8 THEN TO_DATE('121015','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 9 THEN TO_DATE('130114','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 10 THEN TO_DATE('130417','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 11 THEN TO_DATE('130716','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 12 THEN TO_DATE('131015','YYMMDD')
                        END as PSU_DATE
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                        or COMMENTS like '%DATABASE PATCH SET UPDATE%'
                ),
                STATE as
                (
                        select 
                        case   
                                WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
                                ELSE 'OK' 
                        END as "STATUS"
                        from DATA
                )
                select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS 
                into :DESCRIPTION,:PSU_DATE,:STATUS 
                from PSU,DATA,STATE;
         ELSIF ( :EXTENDVERSION = '11.2.0.2.0' AND :EXIST > 0 AND :BP like '%BP%') THEN
         with PSU as
         (
             select max(COMMENTS) as COMMENTS
             from  registry$history
             where ACTION='APPLY'
             and COMMENTS like 'BP%'
        ),
         DATA as
         (
            select
             case
                 WHEN MAX(COMMENTS) = 'BP22' THEN TO_DATE('151013','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP21' THEN TO_DATE('150713','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP20' THEN TO_DATE('150413','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP19' THEN TO_DATE('150113','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP18' THEN TO_DATE('151012','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP17' THEN TO_DATE('150712','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP16' THEN TO_DATE('150412','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP15' THEN TO_DATE('150112','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP14' THEN TO_DATE('150112','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP13' THEN TO_DATE('151011','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP12' THEN TO_DATE('151011','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP11' THEN TO_DATE('150711','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP10' THEN TO_DATE('150711','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP9'  THEN TO_DATE('150711','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP8'  THEN TO_DATE('150411','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP7'  THEN TO_DATE('150411','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP6'  THEN TO_DATE('150411','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP5'  THEN TO_DATE('180311','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP4'  THEN TO_DATE('170211','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP3'  THEN TO_DATE('180111','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP2'  THEN TO_DATE('171210','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP1'  THEN TO_DATE('081210','YYMMDD')
                 ELSE
                    TO_DATE(substr(COMMENTS,3,7), 'YYMMDD')				 
             END as PSU_DATE
             from  registry$history
             where ACTION='APPLY'
             and COMMENTS like 'BP%'
         ),
         STATE as
         (
             select
             case
                 WHEN PSU_DATE < SYSDATE-180 THEN 'KO'
                 ELSE 'OK'
             END as "STATUS"
             from DATA
         )
         select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS
         into :DESCRIPTION,:PSU_DATE,:STATUS
         from PSU,DATA,STATE;								
   ELSIF ( :EXTENDVERSION = '11.2.0.1.0' AND :EXIST > 0  AND :BP not like '%BP%') THEN 
        with PSU as
                (
                        select max(COMMENTS) as COMMENTS
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                ),
                DATA as
                (
                        select 
                        case
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 1 THEN TO_DATE('100413','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 2 THEN TO_DATE('100713','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 3 THEN TO_DATE('101012','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 4 THEN TO_DATE('110118','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 5 THEN TO_DATE('110418','YYMMDD')
                                WHEN substr(MAX(COMMENTS), - instr(MAX(reverse(COMMENTS)), '.') + 1) = 6 THEN TO_DATE('110718','YYMMDD')
                        END as PSU_DATE
                        from  registry$history
                        where ACTION='APPLY' 
                        and COMMENTS like 'PSU%'
                ),
                STATE as
                (
                        select 
                        case   
                                WHEN PSU_DATE < SYSDATE-180 THEN 'KO' 
                                ELSE 'OK' 
                        END as "STATUS"
                        from DATA
                )
                select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS 
                into :DESCRIPTION,:PSU_DATE,:STATUS 
                from PSU,DATA,STATE;
         ELSIF ( :EXTENDVERSION = '11.2.0.1.0' AND :EXIST > 0 AND :BP like '%BP%') THEN
         with PSU as
         (
             select max(COMMENTS) as COMMENTS
             from  registry$history
             where ACTION='APPLY'
             and COMMENTS like 'BP%'
        ),
         DATA as
         (
            select
             case
                 WHEN MAX(COMMENTS) = 'BP12' THEN TO_DATE('150711','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP11' THEN TO_DATE('150711','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP10' THEN TO_DATE('150411','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP9'  THEN TO_DATE('150711','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP8'  THEN TO_DATE('150111','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP7'  THEN TO_DATE('151010','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP6'  THEN TO_DATE('150710','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP5'  THEN TO_DATE('081210','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP4'  THEN TO_DATE('081210','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP3'  THEN TO_DATE('150410','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP2'  THEN TO_DATE('110310','YYMMDD')
                 WHEN MAX(COMMENTS) = 'BP1'  THEN TO_DATE('100210','YYMMDD')
                 ELSE
                    TO_DATE(substr(COMMENTS,3,7), 'YYMMDD')				 
             END as PSU_DATE
             from  registry$history
             where ACTION='APPLY'
             and COMMENTS like 'BP%'
         ),
         STATE as
         (
             select
             case
                 WHEN PSU_DATE < SYSDATE-180 THEN 'KO'
                 ELSE 'OK'
             END as "STATUS"
             from DATA
         )
         select distinct COMMENTS as DESCRIPTION,PSU_DATE,STATUS
         into :DESCRIPTION,:PSU_DATE,:STATUS
         from PSU,DATA,STATE;										
    ELSE
        select 'N/A','N/A','N/A' 
                into :DESCRIPTION,:PSU_DATE,:STATUS 
                from dual;
   END IF; 

END;
/

col Description for a70
col PSU for a40
col STATUS for a40
select :DESCRIPTION as Description 
           ,:PSU_DATE as PSU
--         ,:STATUS as STATUS
from dual WHERE :PSU_DATE != 'N/A';
EXIT