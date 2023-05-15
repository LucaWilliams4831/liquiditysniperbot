package controllers

import (
	"fmt"
	"os"
	"time"
	"math/big"
	"context"
	"encoding/hex"
	"os/exec"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v4"

	"github.com/LucaWilliams4831/uniswap-pancakeswap-tradingbot/liquiditysniperbot/database"
	"github.com/LucaWilliams4831/uniswap-pancakeswap-tradingbot/liquiditysniperbot/models"
)

type Account struct {
	Id        uint   	`json:"id"`
	Address   string 	`json:"address"`
	Status    int 		`json:"status"`
	Fee 	  string    `json:"fee"`  // -: received, <address>: awaiting receipt
	Type 	  string    `json:"type"` // 0: Account, 1: Contract
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var Secretkey = os.Getenv("SECRETKEY")
	
func GetAccounts(c *fiber.Ctx) error {
		
	var accounts []models.Account
	database.DB.Find(&accounts)
	fmt.Println("Found accounts..")
	c.Status(fiber.StatusOK)

	return c.JSON(fiber.Map{
		"message": "Found accounts",
		"data": accounts,
	})
}

func SendFee(c *fiber.Ctx) error {
	paramid := c.Params("id")	

	var account models.Account
	fmt.Println(paramid)
	database.DB.Where("id = ?", paramid).First(&account)
	startString := string(account.Address[0])
	startFeeString := string(account.Fee[0])

	if startString == "d" {
		fmt.Println("Sending Cosmos...")
		fmt.Println(string(account.Address))
		suc := SendCosmos(account.Address);
		
		if (account.Address != account.Fee && startFeeString == "d" ) {
			fmt.Println("---------------")
			fmt.Println(string(account.Fee))
			suc = SendCosmos(account.Fee);
		}

		if suc {
			fmt.Println("Sent Cosmos.")
			account.Fee = "-"
			database.DB.Save(&account)
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"message": "Account Fee added successfully",
			})
		} else {
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"message": "Cosmos Account Fee Failed",
			})
		}
	} else {
		fmt.Println("Sending Eth...")
		suc := SendETH(account.Address);
		
		if (account.Address != account.Fee && startFeeString == "0") {
			fmt.Println("**********************")
			fmt.Println(string(account.Fee))
			suc = SendETH(account.Fee);
		}

		if suc {
			fmt.Println("Sent Eth.")
			account.Fee = "-"
			database.DB.Save(&account)
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"message": "Account Fee added successfully",
			})
		} else {
			c.Status(fiber.StatusOK)
			return c.JSON(fiber.Map{
				"message": "Account Fee Failed",
			})
		}
	}
}

func UpdateAccount(c *fiber.Ctx) error {
	paramid := c.Params("id")
	data := new(Account)

	// cookie := c.Cookies("jwt")

	// token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
	// 	return []byte(Secretkey), nil
	// })

	// if err != nil {
	// 	c.Status(fiber.StatusUnauthorized)
	// 	return c.JSON(fiber.Map{
	// 		"message": "unauthenticated",
	// 	})
	// }

	if err := c.BodyParser(data); err != nil {
		return err
	}
	
	accountstatus := data.Status

	// claims, _ := token.Claims.(*jwt.StandardClaims)
	// user_ID := claims.Issuer

	// var user models.User
	// database.DB.Where("id = ?", user_ID).Find(&user)
	// database.DB.Where("id = ?", data.Id).Find(&user)
	
	// if user.Id == 0 {
	// 	c.Status(fiber.StatusUnauthorized)
	// 	return c.JSON(fiber.Map{
	// 		"message": "unauthenticated",
	// 	})
	// }
	var account models.Account
	database.DB.Where("id = ?", paramid).First(&account)
	account.Status = accountstatus
	account.UpdatedAt = time.Now()
	database.DB.Save(&account)
	c.Status(fiber.StatusOK)
	return c.JSON(fiber.Map{
		"message": "Account Updated successfully",
	})
}

func AddAccount(c *fiber.Ctx) error {
	paramid := c.Params("address")

	var accounts []models.Account
	var newaccount models.Account
	database.DB.Where("address = ?", paramid).Find(&accounts)

	if len(accounts) == 0 {
		newaccount.Address = paramid
		newaccount.Fee = "Not_Approved_yet"
		database.DB.Create(&newaccount)
		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{
			"message": "Account Updated successfully",
			"data": newaccount,
		})
	}
	
	return c.JSON(fiber.Map{
		"message": "Account Duplicated",
	})
}
	
func FilterAccounts(c *fiber.Ctx) error {
	paramkeyword := c.Params("keyword")
		
	var accounts []models.Account
	database.DB.Where("address LIKE ?", "%" + paramkeyword + "%").Find(&accounts)
	c.Status(fiber.StatusOK)

	return c.JSON(fiber.Map{
		"message": "Found accounts",
		"data": accounts,
	})
}

func SendETH(addr string) bool {
	// Connect to Ethereum network

	client, err := ethclient.Dial("http://3.145.87.221:8545")
	if err != nil {
			fmt.Println("Error connecting to Ethereum network:", err)
			return false
	}
	
	// Load wallet  
	privateKeyBytes, _ := hex.DecodeString("DC84C4FC68B190056B9D6D697DA2417089426D2E2F90E3C61417DAE06D48434A")
	privateKey, _ := crypto.ToECDSA(privateKeyBytes)
	walletAddress := common.HexToAddress("0x044204e7E8d4F8F18E3164B7dFC1f8D0Ac550337") // Replace with actual wallet address

	// Get nonce
	nonce, err := client.PendingNonceAt(context.Background(), walletAddress)
	if err != nil {
			fmt.Println("Error getting nonce:", err)
			return false
	}
	
	// Create transaction
	destinationAddress := common.HexToAddress(addr) // Replace with actual destination address
	value := big.NewInt(1000000000000000000) // 1 ETH
	gasLimit := uint64(210000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("Error getting gas price:", err)
		return false
	}
	tx := types.NewTransaction(nonce, destinationAddress, value, gasLimit, gasPrice, nil)
	
	chainID := big.NewInt(1)
	// Sign transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		fmt.Println("Error signing transaction:", err)
		return false
	}

	fmt.Println("************* Now sending transaction *********")
	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
			fmt.Println("Error sending transaction:", err)
			return false
	}

	fmt.Println("Transaction sent:", signedTx.Hash().Hex())
	fmt.Println("-------------------------------------------")

	return true
}

func SendCosmos(addr string) bool {
	cmd := exec.Command("sudo digitaldollard", "tx", "bank", "send", "dd1q3pqfelg6nu0rr33vjmals0c6zk92qehfaxy0c", addr, "1karma", "-y")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(string(out))
	return true
}