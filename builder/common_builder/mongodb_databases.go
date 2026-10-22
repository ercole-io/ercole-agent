// Copyright (c) 2023 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ercole-io/ercole-agent/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/hashicorp/go-multierror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const dbConfig = "config"
const dbLocal = "local"
const dbAdmin = "admin"

func (b *CommonBuilder) getMongoDBFeature(hostname string) (*model.MongoDBFeature, error) {
	date := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s][AGEN][INFO]Fetching MongoDB data...\n", date[2:])

	mongoDB := &model.MongoDBFeature{
		Instances: []model.MongoDBInstance{},
	}

	var merr error

	ctx := context.TODO()

	for i, conf := range b.configuration.Features.MongoDB.Instances {
		client := connection(ctx, conf, b)
		defer func() {
			if err := client.Disconnect(ctx); err != nil {
				b.log.Errorf("Can't disconnect: host %s", hostname)
			}
		}()

		instance := model.MongoDBInstance{}

		var countDBs int

		var commandResultBuildInfo bson.M

		commandbuildInfo := bson.D{bson.E{Key: "buildInfo", Value: 1}}

		var commandResultServerStatus bson.M

		commandServerStatus := bson.D{{Key: "serverStatus", Value: 1}}

		var commandResultDBStats bson.M

		commandDBStats := bson.D{{Key: "dbStats", Value: 1}, {Key: "freeStorage", Value: 1}}

		var commandResultUsersInfo bson.M

		commandUsersInfo := bson.D{{Key: "usersInfo", Value: 1}}

		commandShard := bson.D{{Key: "listShards", Value: 1}}

		commandHelloResult := bsonx.Doc{{Key: "hello", Value: bsonx.Int32(1)}}

		errVersion := client.Database(dbAdmin).RunCommand(ctx, commandbuildInfo).Decode(&commandResultBuildInfo)
		if errVersion != nil {
			b.log.Error(errVersion)
			merr = multierror.Append(merr, errVersion)
		}

		instance.Version = commandResultBuildInfo["version"].(string)

		err := client.Database(dbAdmin).RunCommand(ctx, commandServerStatus).Decode(&commandResultServerStatus)
		if err != nil {
			b.log.Error(err)
			merr = multierror.Append(merr, err)
		}

		statusConnection := model.ServerStatusConnection{}

		conn := commandResultServerStatus["connections"].(primitive.M)
		statusConnection.Current = toInt32(conn["current"])
		statusConnection.Available = toInt32(conn["available"])
		statusConnection.TotalCreated = toInt32(conn["totalCreated"])

		if v, ok := conn["active"]; ok {
			statusConnection.Active = toInt32(v)
		}

		if v, ok := conn["threaded"]; ok {
			statusConnection.Threaded = toInt32(v)
		}

		if v, ok := conn["exhaustIsMaster"]; ok {
			statusConnection.ExhaustIsMaster = toInt32(v)
		}

		if v, ok := conn["exhaustHello"]; ok {
			statusConnection.ExhaustHello = toInt32(v)
		}

		if v, ok := conn["awaitingTopologyChanges"]; ok {
			statusConnection.AwaitingTopologyChanges = toInt32(v)
		}

		if conn["loadBalanced"] != nil {
			statusConnection.LoadBalanced = toInt32(conn["loadBalanced"])
		}

		instance.StatusConnection = statusConnection

		errHello := client.Database(dbAdmin).RunCommand(ctx, commandHelloResult).Decode(&instance.ReplicaSet)
		if errHello != nil {
			b.log.Error(errHello)
			merr = multierror.Append(merr, errHello)
		}

		errShard := client.Database(dbAdmin).RunCommand(ctx, commandShard).Decode(&instance.ShardList)
		if errShard != nil && errShard.Error() != "(CommandNotFound) no such command: 'listShards'" {
			b.log.Error(errShard)
			merr = multierror.Append(merr, errShard)
		}

		filter := bson.D{{}}
		dbs, errList := client.ListDatabaseNames(ctx, filter)

		if errList != nil {
			b.log.Error(errList)
			merr = multierror.Append(merr, errList)
		}

		for _, db := range dbs {
			if db != dbAdmin && db != dbConfig && db != dbLocal {
				dbStats := model.DBStats{}
				dbStats.DBName = db
				dbStats.Charset = "UTF8"

				countDBs += 1

				errDB := client.Database(db).RunCommand(ctx, commandDBStats).Decode(&commandResultDBStats)
				if errDB != nil {
					b.log.Error(errDB)
					merr = multierror.Append(merr, errDB)
				}

				if commandResultDBStats["raw"] != nil {
					raw := commandResultDBStats["raw"].(primitive.M)
					for k := range raw {
						rawContent := raw[k].(primitive.M)
						dbStats.Collections = toInt32(rawContent["collections"])
						dbStats.Views = toInt32(rawContent["views"])
						dbStats.FsUsedSize = toFloat64(rawContent["fsUsedSize"])
						dbStats.FsTotalSize = toFloat64(rawContent["fsTotalSize"])
						dbStats.FreeStorageSize = toFloat64(rawContent["freeStorageSize"])
						dbStats.IndexFreeStorageSize = toFloat64(rawContent["indexFreeStorageSize"])
						dbStats.TotalFreeStorageSize = toFloat64(rawContent["totalFreeStorageSize"])
						dbStats.DataSize = toFloat64(commandResultDBStats["dataSize"])
						dbStats.IndexSize = toFloat64(commandResultDBStats["indexSize"])
						dbStats.StorageSize = toFloat64(commandResultDBStats["storageSize"])
						dbStats.TotalSize = toFloat64(commandResultDBStats["totalSize"])
					}
				} else {
					dbStats.Collections = toInt32(commandResultDBStats["collections"])
					dbStats.Views = toInt32(commandResultDBStats["views"])
					dbStats.FsUsedSize = toFloat64(commandResultDBStats["fsUsedSize"])
					dbStats.FsTotalSize = toFloat64(commandResultDBStats["fsTotalSize"])
					dbStats.FreeStorageSize = toFloat64(commandResultDBStats["freeStorageSize"])
					dbStats.IndexFreeStorageSize = toFloat64(commandResultDBStats["indexFreeStorageSize"])
					dbStats.TotalFreeStorageSize = toFloat64(commandResultDBStats["totalFreeStorageSize"])
					dbStats.DataSize = toFloat64(commandResultDBStats["dataSize"])
					dbStats.IndexSize = toFloat64(commandResultDBStats["indexSize"])
					dbStats.StorageSize = toFloat64(commandResultDBStats["storageSize"])
					dbStats.TotalSize = toFloat64(commandResultDBStats["totalSize"])
				}

				dbStats.Objects = toInt32(commandResultDBStats["objects"])
				dbStats.Indexes = toInt32(commandResultDBStats["indexes"])

				errUserInfo := client.Database(db).RunCommand(ctx, commandUsersInfo).Decode(&commandResultUsersInfo)
				if errUserInfo != nil {
					b.log.Error(errUserInfo)
					merr = multierror.Append(merr, errUserInfo)
				}

				dbStats.Users = len(commandResultUsersInfo["users"].(primitive.A))

				shardDBs, errShardDbs := client.Database(dbConfig).Collection("databases").Find(ctx, bson.M{"_id": db})
				if errShardDbs != nil && errShardDbs != mongo.ErrNoDocuments {
					b.log.Error(errShardDbs)
					merr = multierror.Append(merr, errShardDbs)
				}

				if errDBs := shardDBs.All(ctx, &dbStats.ShardDBs); errDBs != nil {
					b.log.Error(errDBs)
					merr = multierror.Append(merr, errDBs)
				}

				instance.Stats = append(instance.Stats, dbStats)
			}
		}

		instance.Dbs = countDBs
		instance.Name = b.configuration.Features.MongoDB.Instances[i].Host + ":" + b.configuration.Features.MongoDB.Instances[i].Port

		mongoDB.Instances = append(mongoDB.Instances, instance)
	}

	return mongoDB, merr
}

func connection(ctx context.Context, conf config.MongoDBInstanceConnection, b *CommonBuilder) *mongo.Client {
	var connectionString string

	mongoDBUser := conf.User
	mongoDBPWD := conf.Password
	mongoDBHost := conf.Host
	mongoDBPort := conf.Port
	mongoDBDirectConn := conf.DirectConnection

	if mongoDBUser != "" && mongoDBPWD != "" {
		connectionString = fmt.Sprintf("mongodb://%s:%s@%s:%s/?directConnection=%s", mongoDBUser, mongoDBPWD, mongoDBHost, mongoDBPort, strconv.FormatBool(mongoDBDirectConn))
	} else {
		connectionString = fmt.Sprintf("mongodb://%s:%s/?directConnection=%s", mongoDBHost, mongoDBPort, strconv.FormatBool(mongoDBDirectConn))
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		b.log.Fatalf("Can't connect to MongoDB: host %s port %s", mongoDBHost, mongoDBPort)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		b.log.Fatalf("Can't ping readpref.Primary: host %s port %s", mongoDBHost, mongoDBPort)
	}

	return client
}

func toInt32(v any) int32 {
	switch n := v.(type) {
	case int32:
		return n
	case int64:
		return int32(n)
	case float64:
		return int32(n)
	case primitive.Decimal128:
		f, err := strconv.ParseFloat(n.String(), 64)
		if err != nil {
			return 0
		}

		return int32(f)
	default:
		return 0
	}
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	case primitive.Decimal128:
		f, err := strconv.ParseFloat(n.String(), 64)
		if err != nil {
			return 0
		}

		return f
	default:
		return 0
	}
}
