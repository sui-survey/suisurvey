package api

import (
	"context"
	"fmt"
	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
	"kevinsheeran/walrus-backend/model"
)

// blobId string, name string, expiration string, minParticipants string, maxParticipants string, reward string, state string, clock string, contractInteraction string
func createSurvey(data *model.CreateFormDto, name string, blobId string) {
	var ctx = context.Background()
	var cli = sui.NewSuiClient(constant.SuiTestnetEndpoint)

	signerAccount, err := signer.NewSignertWithMnemonic("fire quarter pelican often evidence toss merge feel remember absurd learn glove")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	priKey := signerAccount.PriKey
	fmt.Printf("signerAccount.Address: %s\n", signerAccount.Address)

	gasObj := "0x55ccc9b3cda7a3d5d84e6282b230d5e72231ea244758d49e06018fbc5034b294"

	rsp, err := cli.MoveCall(ctx, models.MoveCallRequest{
		Signer:          signerAccount.Address,
		PackageObjectId: "0xe6ec46bacdf52039961a503a65c838987b54b0fe1b5d27b6bd9991c2e5df3fb7",
		Module:          "suisurvey",
		Function:        "create_survey",
		TypeArguments:   []interface{}{},
		Arguments: []interface{}{
			data.Id,     // Creator
			name,        // Name
			"864000000", // Replace nil with 0 or a default value if U64 is expected
			"1",
			"10",   // Replace nil with 0 or a default value if U64 is expected
			blobId, // blobId
			"0xe6ec46bacdf52039961a503a65c838987b54b0fe1b5d27b6bd9991c2e5df3fb7", // Replace nil with 0 or a default value if U64 is expected
			"0", // Replace nil with 0 or a default value if U64 is expected
			"0x56f99f6bddabda730c57fe729d6ff7586093b01e00de876a1766f3da0108ec45",
			"0x6",
		},
		Gas:       gasObj,
		GasBudget: "100000000",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// see the successful transaction url: https://explorer.sui.io/txblock/CD5hFB4bWFThhb6FtvKq3xAxRri72vsYLJAVd7p9t2sR?network=testnet
	rsp2, err := cli.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: rsp,
		PriKey:      priKey,
		// only fetch the effects field
		Options: models.SuiTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForLocalExecution",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utils.PrettyPrint(rsp2)
}

func participateSurvey() {
	var ctx = context.Background()
	var cli = sui.NewSuiClient(constant.SuiTestnetEndpoint)

	signerAccount, err := signer.NewSignertWithMnemonic("fire quarter pelican often evidence toss merge feel remember absurd learn glove")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	priKey := signerAccount.PriKey
	fmt.Printf("signerAccount.Address: %s\n", signerAccount.Address)

	gasObj := "0x55ccc9b3cda7a3d5d84e6282b230d5e72231ea244758d49e06018fbc5034b294"

	rsp, err := cli.MoveCall(ctx, models.MoveCallRequest{
		Signer:          signerAccount.Address,
		PackageObjectId: "0xe6ec46bacdf52039961a503a65c838987b54b0fe1b5d27b6bd9991c2e5df3fb7",
		Module:          "suisurvey",
		Function:        "participate_survey",
		TypeArguments:   []interface{}{},
		Arguments: []interface{}{
			"0xb6ba387da7991e22ef89e0a2d0323f7102990bfb69157f983ef9aa8326d793e1", // participant address
			"test",      // Form object address
			"864000000", // Creator address
			"0K7wNVHRw8_FHrJOh8ItI94EPpJiuTldRj3mke4vrNGY",                       // blobId
			"0x56f99f6bddabda730c57fe729d6ff7586093b01e00de876a1766f3da0108ec45", //
			"0x6", // state
		},
		Gas:       gasObj,
		GasBudget: "100000000",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// see the successful transaction url: https://explorer.sui.io/txblock/CD5hFB4bWFThhb6FtvKq3xAxRri72vsYLJAVd7p9t2sR?network=testnet
	rsp2, err := cli.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: rsp,
		PriKey:      priKey,
		// only fetch the effects field
		Options: models.SuiTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForLocalExecution",
	})

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utils.PrettyPrint(rsp2)
}
