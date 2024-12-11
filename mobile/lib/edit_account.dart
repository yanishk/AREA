import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'package:namer_app/login_page.dart';
import 'package:fluttertoast/fluttertoast.dart';

import "custom_widget.dart";
import "footer_page.dart";
import "editPassword.dart";
import "globalData.dart";

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
  final String apiUrl = dotenv.get('SERVER_LINK',
      fallback: 'Cannot find .env variable "SERVER_LINK"');
  String apiRoute = 'user/delete-account';
  bool isInit = true;
  final TextEditingController emailController = new TextEditingController(
    text: "",
  );

  Future<void> _deleteAccount() async {
    developer.log('SERVER_LINK : $apiUrl');
    try {
      developer.log('SERVER_LINK + redirection : ${apiUrl + apiRoute}');
      final response = await http.delete(
        Uri.parse(apiUrl + apiRoute),
        headers: <String, String>{
          'Content-Type': 'application/json; charset=UTF-8',
          'Cookie': Globaldata.JWToken,
        },
      );

      if (response.statusCode == 200) {
        developer
            .log("Account deleted by user, JWToken : ${Globaldata.JWToken}");
        setState(() {
          Globaldata.JWToken = "";
          Fluttertoast.showToast(
              msg: "Account deleted",
              toastLength: Toast.LENGTH_SHORT,
              timeInSecForIosWeb: 1,
              backgroundColor: Colors.red,
              textColor: Colors.white,
              fontSize: 24.0);

          Navigator.push(
            Globaldata.myContext,
            MaterialPageRoute(builder: (context) => LoginPage()),
          );
        });
      } else {
        final responseData = jsonDecode(response.body);
        setState(() {
          Fluttertoast.showToast(
              msg: '${responseData['error']}',
              toastLength: Toast.LENGTH_SHORT,
              timeInSecForIosWeb: 1,
              backgroundColor: Colors.red,
              textColor: Colors.white,
              fontSize: 24.0);
        });
      }
    } catch (e) {
      setState(() {
        developer.log('$e');
        if ('$e'.split(' ').contains("SocketException:") == true ||
            ('$e'.split(' ').contains("Invalid") == true &&
                '$e'.split(' ').contains("argument(s):") == true)) {
          Fluttertoast.showToast(
              msg: Globaldata.ErrorServerConnection,
              toastLength: Toast.LENGTH_SHORT,
              timeInSecForIosWeb: 1,
              backgroundColor: Colors.red,
              textColor: Colors.white,
              fontSize: 24.0);
        } else {
          Fluttertoast.showToast(
              msg: '$e',
              toastLength: Toast.LENGTH_SHORT,
              timeInSecForIosWeb: 1,
              backgroundColor: Colors.red,
              textColor: Colors.white,
              fontSize: 24.0);
        }
      });
    }
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
            _deleteAccount();
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
            Globaldata.JWToken = "";
            Navigator.of(context).pop();
            Navigator.push(
              Globaldata.myContext,
              MaterialPageRoute(builder: (context) => LoginPage()),
            );
          },
          secondChoiceonPressed: () {
            Navigator.of(context).pop();
          },
        );
      },
    );
  }

  Future<void> getEmail() async {
    final response = await GetApi.get("user/email");
    emailController.text = response["email"];
  }

  @override
  Widget build(BuildContext context) {
    if (isInit == true) {
      isInit = false;
      getEmail();
    }
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            AutoSizeText(
              "Account",
              style: TextStyle(
                fontSize: 30,
                fontWeight: FontWeight.bold,
                color: Colors.black,
              ),
              maxLines: 1,
            ),
            SizedBox(width: 16),
            MyDivider(),
            Padding(padding: EdgeInsets.only(top: 30)),
            LeftWidget(
              child: AutoSizeText(
                "Email",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
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
            LeftWidget(
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
            ),
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
      bottomNavigationBar: FooterNavigationBar(currentIndex: 4),
    );
  }
}
