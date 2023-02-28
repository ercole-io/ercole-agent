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

	feature := model.MongoDBFeature{}

	var merr error

	ctx := context.TODO()

	client := connection(ctx, b)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			b.log.Errorf("Can't disconnect: host %s", hostname)
		}
	}()

	var countDBs int

	commandServerStatus := bson.D{{Key: "serverStatus", Value: 1}}

	var commandResultDBStats bson.M

	commandDBStats := bson.D{{Key: "dbStats", Value: 1}, {Key: "freeStorage", Value: 1}}

	var commandResultUsersInfo bson.M

	commandUsersInfo := bson.D{{Key: "usersInfo", Value: 1}}

	commandShard := bson.D{{Key: "listShards", Value: 1}}

	commandHelloResult := bsonx.Doc{{Key: "hello", Value: bsonx.Int32(1)}}

	err := client.Database(dbAdmin).RunCommand(ctx, commandServerStatus).Decode(&feature.Status)
	if err != nil {
		b.log.Error(err)
		merr = multierror.Append(merr, err)
	}

	errHello := client.Database(dbAdmin).RunCommand(ctx, commandHelloResult).Decode(&feature.ReplicaSet)
	if errHello != nil {
		b.log.Error(errHello)
		merr = multierror.Append(merr, errHello)
	}

	errShard := client.Database(dbAdmin).RunCommand(ctx, commandShard).Decode(&feature.ShardList)
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
				fmt.Println(errDB)
			}

			if commandResultDBStats["raw"] != nil {
				raw := commandResultDBStats["raw"].(primitive.M)
				for k := range raw {
					rawContent := raw[k].(primitive.M)
					dbStats.Collections = rawContent["collections"].(int32)
					dbStats.Views = rawContent["views"].(int32)
					dbStats.FsUsedSize = rawContent["fsUsedSize"].(float64)
					dbStats.FsTotalSize = rawContent["fsTotalSize"].(float64)
					dbStats.FreeStorageSize = rawContent["freeStorageSize"].(float64)
					dbStats.IndexFreeStorageSize = rawContent["indexFreeStorageSize"].(float64)
					dbStats.TotalFreeStorageSize = rawContent["totalFreeStorageSize"].(float64)
					dbStats.DataSize = float64(commandResultDBStats["dataSize"].(int32))
					dbStats.IndexSize = float64(commandResultDBStats["indexSize"].(int32))
					dbStats.StorageSize = float64(commandResultDBStats["storageSize"].(int32))
					dbStats.TotalSize = float64(commandResultDBStats["totalSize"].(int32))
				}
			} else {
				dbStats.Collections = commandResultDBStats["collections"].(int32)
				dbStats.Views = commandResultDBStats["views"].(int32)
				dbStats.FsUsedSize = commandResultDBStats["fsUsedSize"].(float64)
				dbStats.FsTotalSize = commandResultDBStats["fsTotalSize"].(float64)
				dbStats.FreeStorageSize = commandResultDBStats["freeStorageSize"].(float64)
				dbStats.IndexFreeStorageSize = commandResultDBStats["indexFreeStorageSize"].(float64)
				dbStats.TotalFreeStorageSize = commandResultDBStats["totalFreeStorageSize"].(float64)
				dbStats.DataSize = commandResultDBStats["dataSize"].(float64)
				dbStats.IndexSize = commandResultDBStats["indexSize"].(float64)
				dbStats.StorageSize = commandResultDBStats["storageSize"].(float64)
				dbStats.TotalSize = commandResultDBStats["totalSize"].(float64)
			}

			dbStats.Objects = commandResultDBStats["objects"].(int32)
			dbStats.Indexes = commandResultDBStats["indexes"].(int32)

			errUserInfo := client.Database(db).RunCommand(ctx, commandUsersInfo).Decode(&commandResultUsersInfo)
			if errUserInfo != nil {
				fmt.Println(errUserInfo)
			}

			dbStats.Users = len(commandResultUsersInfo["users"].(primitive.A))

			shardDBs, errShardDbs := client.Database(dbConfig).Collection("databases").Find(ctx, bson.M{"_id": db})
			if errShardDbs != nil && errShardDbs != mongo.ErrNoDocuments {
				fmt.Println(errShardDbs)
			}

			if errDBs := shardDBs.All(ctx, &dbStats.ShardDBs); err != nil {
				fmt.Println(errDBs)
			}

			feature.Stats = append(feature.Stats, dbStats)
		}
	}

	feature.Status.Dbs = countDBs

	return &feature, merr
}

func connection(ctx context.Context, b *CommonBuilder) *mongo.Client {
	var connectionString string

	mongoDBUser := b.configuration.Features.MongoDB.User
	mongoDBPWD := b.configuration.Features.MongoDB.Password
	mongoDBHost := b.configuration.Features.MongoDB.Host
	mongoDBPort := b.configuration.Features.MongoDB.Port
	mongoDBDirectConn := b.configuration.Features.MongoDB.DirectConnection

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
