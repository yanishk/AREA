import 'package:flutter/material.dart';
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:auto_size_text/auto_size_text.dart';
import 'package:fluttertoast/fluttertoast.dart';
import "package:http/src/response.dart";
import 'package:flutter_phoenix/flutter_phoenix.dart';

import "../custom_widget.dart";
import "../footer_page.dart";
import "editPassword.dart";
import "../globalData.dart";
import "../apiCall/apiRequest.dart";

class EditAccountPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  EditAccount createState() => EditAccount();
}

class EditAccount extends State<APICall> {
  String apiRoute = 'user';
  final TextEditingController emailController = new TextEditingController(
    text: Globaldata.userInfos.email,
  );

  displayToastError(String text) {
    Fluttertoast.showToast(
      msg: text,
      toastLength: Toast.LENGTH_SHORT,
      timeInSecForIosWeb: 1,
      backgroundColor: Colors.red,
      textColor: Colors.white,
      fontSize: 24.0
    );
  }

  Future<void> deleteAccount() async {
    Response response = await ApiRequest.delete(apiRoute);
    if (response.statusCode == 504) {
      displayToastError(response.body);
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      displayToastError('${responseData['error']}');
      developer.log(
          "deleteAccount FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return;
    }
    setState(() {
      Globaldata.JWToken = "";
      Fluttertoast.showToast(
        msg: "Account deleted",
        toastLength: Toast.LENGTH_SHORT,
        timeInSecForIosWeb: 1,
        backgroundColor: Colors.red,
        textColor: Colors.white,
        fontSize: 24.0
      );
    });
    Phoenix.rebirth(Globaldata.myContext);
  }

  Future<void> deleteAccountConfirmDialog() async {
    return showDialog<void>(
      context: Globaldata.myContext,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return TwoChoiceAlertDialog(
          firstChoiceText: Text(
            "Delete my account",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white),
          ),
          secondChoiceText: Text(
            "Cancel",
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.black,
            ),
          ),
          firstChoiceonPressed: () {
            deleteAccount();
            Navigator.of(context).pop();
          },
          secondChoiceonPressed: () {
            Navigator.of(context).pop();
          },
        );
      },
    );
  }

  Future<void> signOutConfirmDialog() async {
    return showDialog<void>(
      context: Globaldata.myContext,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return TwoChoiceAlertDialog(
          firstChoiceText: Text(
            "Sign out",
            style: TextStyle(
                fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white),
          ),
          secondChoiceText: Text(
            "Cancel",
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.black,
            ),
          ),
          firstChoiceonPressed: () {
            Globaldata.resetActionReactionServicesList();
            Globaldata.JWToken = "";
            Navigator.of(context).pop();
            Phoenix.rebirth(Globaldata.myContext);
          },
          secondChoiceonPressed: () {
            Navigator.of(context).pop();
          },
        );
      },
    );
  }

  Widget displayChangePasswordButton() {
    if (Globaldata.userInfos.connectiontype != "basic") {
      return Column();
    }
    return LeftWidget(
      child: CustomTextButton(
        text: Text(
          "Change password",
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Color.fromARGB(255, 0, 160, 255),
          ),
        ),
        onPressed: () {
          Navigator.push(
            Globaldata.myContext,
            MaterialPageRoute(builder: (context) => EditPasswordPage()),
          );
        },
        padding: EdgeInsets.zero,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            AutoSizeText("Account settings",
              style: TextStyle(fontSize: 30,fontWeight: FontWeight.w900,color: Colors.black,),
              maxLines: 1,
            ),
            SizedBox(width: 16),
            MyDivider(),
            Padding(padding: EdgeInsets.only(top: 30)),
            LeftWidget(
              child: AutoSizeText(
                "Email",
                style: TextStyle(fontSize: 25,fontWeight: FontWeight.bold,color: Colors.black,),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
              hint: 'Email',
              padding: EdgeInsets.symmetric(horizontal: 15),
              controller: emailController,
            ),
            Padding(padding: EdgeInsets.only(top: 20)),
            LeftWidget(
              child: AutoSizeText(
                "Password",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
              hint: '••••••••',
              obscureText: true,
              enabled: false,
              fillCollor: Color.fromARGB(255, 238, 238, 238),
              padding: EdgeInsets.symmetric(horizontal: 15),
            ),
            displayChangePasswordButton(),
            Padding(padding: EdgeInsets.only(top: 40)),
            CustomButton(
              backgroundColor: Color.fromARGB(255, 240, 240, 240),
              text: AutoSizeText(
                "Delete my account",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: const Color.fromARGB(255, 255, 0, 0)),
              ),
              onPressed: deleteAccountConfirmDialog,
              size: Size(double.infinity, 55),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
            Padding(padding: EdgeInsets.only(top: 50)),
            CustomButton(
              backgroundColor: Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText(
                "Sign out",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: Colors.white),
              ),
              onPressed: signOutConfirmDialog,
              size: Size(double.infinity, 55),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 3),
    );
  }
}
