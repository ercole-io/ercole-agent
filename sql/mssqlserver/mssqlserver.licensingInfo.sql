DECLARE @temp as TABLE( 
	LogDate datetime, 
	ProcessInfo varchar(100), 
	TextData varchar(max) 
);
DECLARE @ProductCode AS nvarchar(4000),
@EditionType AS nvarchar(4000);

IF  (convert(VARCHAR(2),SERVERPROPERTY('ProductVersion')) >='11')
BEGIN
INSERT INTO @temp
	EXEC master.dbo.sp_readerrorlog @p1 = 0
		,@p2 = 1
		,@p3 = N'licensing'
END
EXEC master.dbo.xp_instance_regread
    N'HKEY_LOCAL_MACHINE',
    N'Software\Microsoft\MSSQLServer\Setup',
    N'ProductCode', 
    @ProductCode output

EXEC master.dbo.xp_instance_regread
    N'HKEY_LOCAL_MACHINE',
    N'Software\Microsoft\MSSQLServer\Setup',
    N'EditionType', 
    @EditionType output
	
SELECT
	SERVERPROPERTY('ProductVersion') AS [ProductVersion],
	@EditionType As [EditionType], 
	@ProductCode As [ProductCode],
	(SELECT top 1
		textData 
	FROM
		@temp
	WHERE 
		textData like 'SQL Server detected%') AS [LicensingInfo]