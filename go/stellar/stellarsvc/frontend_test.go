package stellarsvc

import (
	"context"
	"testing"

	"github.com/keybase/client/go/protocol/stellar1"
	"github.com/keybase/client/go/stellar"
	"github.com/keybase/client/go/stellar/remote"
	"github.com/stretchr/testify/require"
)

func TestGetWalletAccountsLocal(t *testing.T) {
	tcs, cleanup := setupNTests(t, 1)
	defer cleanup()

	_, err := stellar.CreateWallet(context.Background(), tcs[0].G)
	require.NoError(t, err)

	accountID := tcs[0].Backend.AddAccount()

	argImport := stellar1.ImportSecretKeyLocalArg{
		SecretKey:   tcs[0].Backend.SecretKey(accountID),
		MakePrimary: true,
	}
	err = tcs[0].Srv.ImportSecretKeyLocal(context.Background(), argImport)
	require.NoError(t, err)

	tcs[0].Backend.ImportAccountsForUser(tcs[0])

	accts, err := tcs[0].Srv.GetWalletAccountsLocal(context.Background(), 0)
	require.NoError(t, err)

	require.Len(t, accts, 2)
	require.Equal(t, accountID, accts[0].AccountID, accountID)
	require.True(t, accts[0].IsDefault)
	require.Equal(t, "", accts[0].Name) // TODO: once we can set the name on an account, check this
	require.Equal(t, "10,000 XLM", accts[0].BalanceDescription)
	require.False(t, accts[1].IsDefault)
	require.Equal(t, "", accts[1].Name)
	require.Equal(t, "0 XLM", accts[1].BalanceDescription)
}

func TestGetAccountAssetsLocalWithBalance(t *testing.T) {
	tcs, cleanup := setupNTests(t, 1)
	defer cleanup()

	_, err := stellar.CreateWallet(context.Background(), tcs[0].G)
	require.NoError(t, err)

	accountID := tcs[0].Backend.AddAccount()

	argImport := stellar1.ImportSecretKeyLocalArg{
		SecretKey:   tcs[0].Backend.SecretKey(accountID),
		MakePrimary: true,
	}
	err = tcs[0].Srv.ImportSecretKeyLocal(context.Background(), argImport)
	require.NoError(t, err)

	tcs[0].Backend.ImportAccountsForUser(tcs[0])

	assets, err := tcs[0].Srv.GetAccountAssetsLocal(context.Background(), stellar1.GetAccountAssetsLocalArg{AccountID: accountID})
	require.NoError(t, err)

	require.Len(t, assets, 1)
	require.Equal(t, "Lumens", assets[0].Name)
	require.Equal(t, "XLM", assets[0].AssetCode)
	require.Equal(t, "10,000", assets[0].BalanceTotal)
	require.Equal(t, "9,999", assets[0].BalanceAvailableToSend)
	require.Equal(t, "Stellar", assets[0].Issuer)
	require.Equal(t, "USD", assets[0].WorthCurrency)
	require.Equal(t, "$3,183.28", assets[0].Worth)
}

func TestGetAccountAssetsLocalEmptyBalance(t *testing.T) {
	tcs, cleanup := setupNTests(t, 1)
	defer cleanup()

	_, err := stellar.CreateWallet(context.Background(), tcs[0].G)
	require.NoError(t, err)

	accountID := tcs[0].Backend.AddAccountEmpty(t)

	argImport := stellar1.ImportSecretKeyLocalArg{
		SecretKey:   tcs[0].Backend.SecretKey(accountID),
		MakePrimary: true,
	}
	err = tcs[0].Srv.ImportSecretKeyLocal(context.Background(), argImport)
	require.NoError(t, err)

	tcs[0].Backend.ImportAccountsForUser(tcs[0])

	assets, err := tcs[0].Srv.GetAccountAssetsLocal(context.Background(), stellar1.GetAccountAssetsLocalArg{AccountID: accountID})
	require.NoError(t, err)

	require.Len(t, assets, 1)
	require.Equal(t, "Lumens", assets[0].Name)
	require.Equal(t, "XLM", assets[0].AssetCode)
	require.Equal(t, "0", assets[0].BalanceTotal)
	require.Equal(t, "0", assets[0].BalanceAvailableToSend)
	require.Equal(t, "Stellar", assets[0].Issuer)
	require.Equal(t, "USD", assets[0].WorthCurrency)
	require.Equal(t, "$0.00", assets[0].Worth)
}

func TestGetDisplayCurrenciesLocal(t *testing.T) {
	tcs, cleanup := setupNTests(t, 1)
	defer cleanup()

	currencies, err := tcs[0].Srv.GetDisplayCurrenciesLocal(context.Background(), 0)
	require.NoError(t, err)

	require.Len(t, currencies, 32)
	require.Equal(t, "USD ($)", currencies[0].Description)
	require.Equal(t, "USD", currencies[0].Code)
	require.Equal(t, "$", currencies[0].Symbol)
	require.Equal(t, "AUD ($)", currencies[1].Description)
	require.Equal(t, "AUD", currencies[1].Code)
	require.Equal(t, "$", currencies[1].Symbol)
}

func TestChangeWalletName(t *testing.T) {
	tcs, cleanup := setupNTests(t, 1)
	defer cleanup()

	stellar.ServiceInit(tcs[0].G)

	_, err := stellar.CreateWallet(context.Background(), tcs[0].G)
	require.NoError(t, err)

	tcs[0].Backend.ImportAccountsForUser(tcs[0])

	accs, err := tcs[0].Srv.WalletGetAccountsCLILocal(context.Background())
	require.NoError(t, err)
	require.Len(t, accs, 1)
	require.Equal(t, accs[0].Name, "")

	err = tcs[0].Srv.ChangeWalletAccountNameLocal(context.Background(), stellar1.ChangeWalletAccountNameLocalArg{
		AccountID: accs[0].AccountID,
		NewName:   "office lunch money",
	})
	require.NoError(t, err)

	accs, err = tcs[0].Srv.WalletGetAccountsCLILocal(context.Background())
	require.NoError(t, err)
	require.Len(t, accs, 1)
	require.Equal(t, accs[0].Name, "office lunch money")

	// Try invalid argument
	invalidAccID, _ := randomStellarKeypair()
	err = tcs[0].Srv.ChangeWalletAccountNameLocal(context.Background(), stellar1.ChangeWalletAccountNameLocalArg{
		AccountID: invalidAccID,
		NewName:   "savings",
	})
	require.Error(t, err)
}

func TestSetAccountAsDefault(t *testing.T) {
	tcs, cleanup := setupNTests(t, 2)
	defer cleanup()

	_, err := stellar.CreateWallet(context.Background(), tcs[0].G)
	require.NoError(t, err)

	tcs[0].Backend.ImportAccountsForUser(tcs[0])

	accs, err := tcs[0].Srv.WalletGetAccountsCLILocal(context.Background())
	require.NoError(t, err)
	require.Len(t, accs, 1)

	require.True(t, accs[0].IsPrimary)

	// Should work for accounts that are already primary and not post
	// a bundle.
	err = tcs[0].Srv.SetWalletAccountAsDefaultLocal(context.Background(), stellar1.SetWalletAccountAsDefaultLocalArg{
		AccountID: accs[0].AccountID,
	})
	require.NoError(t, err)

	bundle, _, err := remote.Fetch(context.Background(), tcs[0].G)
	require.NoError(t, err)
	require.EqualValues(t, 1, bundle.Revision)

	// Test invalid arguments
	invalidAccID, _ := randomStellarKeypair()
	err = tcs[0].Srv.SetWalletAccountAsDefaultLocal(context.Background(), stellar1.SetWalletAccountAsDefaultLocalArg{
		AccountID: invalidAccID,
	})
	require.Error(t, err)

	err = tcs[0].Srv.SetWalletAccountAsDefaultLocal(context.Background(), stellar1.SetWalletAccountAsDefaultLocalArg{
		AccountID: stellar1.AccountID(""),
	})
	require.Error(t, err)

	additionalAccs := []stellar1.AccountID{
		tcs[0].Backend.AddAccountEmpty(t),
		tcs[0].Backend.AddAccountEmpty(t),
	}

	for _, v := range additionalAccs {
		arg := stellar1.ImportSecretKeyLocalArg{
			SecretKey:   tcs[0].Backend.SecretKey(v),
			MakePrimary: false,
		}
		err = tcs[0].Srv.ImportSecretKeyLocal(context.Background(), arg)
		require.NoError(t, err)
	}

	for i := len(additionalAccs) - 1; i >= 0; i-- {
		v := additionalAccs[i]
		arg := stellar1.SetWalletAccountAsDefaultLocalArg{
			AccountID: v,
		}
		err := tcs[0].Srv.SetWalletAccountAsDefaultLocal(context.Background(), arg)
		require.NoError(t, err)

		accs, err := tcs[0].Srv.WalletGetAccountsCLILocal(context.Background())
		require.NoError(t, err)
		require.Len(t, accs, 3)
		for _, acc := range accs {
			require.Equal(t, acc.IsPrimary, acc.AccountID == v)
		}
	}

	// Expecting additionalAccs[0] to be default account. Lookup
	// public Stellar address as another user.
	u0, err := tcs[1].G.LoadUserByUID(tcs[0].Fu.User.GetUID())
	require.NoError(t, err)
	u0addr := u0.StellarWalletAddress()
	require.NotNil(t, u0addr)
	require.Equal(t, additionalAccs[0], *u0addr)
}

func TestLinkNewWalletAccountLocal(t *testing.T) {
	tcs, cleanup := setupNTests(t, 1)
	defer cleanup()

	// make a wallet
	_, err := stellar.CreateWallet(context.Background(), tcs[0].G)
	require.NoError(t, err)

	// link a new account
	a1, s1 := randomStellarKeypair()
	arg := stellar1.LinkNewWalletAccountLocalArg{
		SecretKey: s1,
		Name:      "my other account",
	}
	linkID, err := tcs[0].Srv.LinkNewWalletAccountLocal(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, a1, linkID)

	tcs[0].Backend.ImportAccountsForUser(tcs[0])

	accts, err := tcs[0].Srv.GetWalletAccountsLocal(context.Background(), 0)
	require.NoError(t, err)

	require.Len(t, accts, 2)
	require.True(t, accts[0].IsDefault)
	require.Equal(t, "", accts[0].Name)
	require.Equal(t, "0 XLM", accts[0].BalanceDescription)
	require.False(t, accts[1].IsDefault)
	require.Equal(t, a1, accts[1].AccountID)
	require.Equal(t, "my other account", accts[1].Name)
	require.Equal(t, "0 XLM", accts[1].BalanceDescription)
}
