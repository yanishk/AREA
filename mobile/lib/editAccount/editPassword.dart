import 'package:flutter/material.dart';
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:auto_size_text/auto_size_text.dart';
import "package:http/src/response.dart";
import 'package:fluttertoast/fluttertoast.dart';

import "../custom_widget.dart";
import "../globalData.dart";
import '../apiCall/apiRequest.dart';

class EditPasswordPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      title: "Change password",
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  EditPassword createState() => EditPassword();
}

class EditPassword extends State<APICall> {
  final TextEditingController oldPasswordController = TextEditingController();
  final TextEditingController newPasswordController = TextEditingController();
  final TextEditingController repeatNewPasswordController = TextEditingController();
  String result = '';
  String notSamePasswd = "";
  String apiRoute = 'user/modify-password';

  @override
  void dispose() {
    oldPasswordController.dispose();
    newPasswordController.dispose();
    super.dispose();
  }

  bool checkSamePassword() {
    if (repeatNewPasswordController.text != newPasswordController.text) {
      setState(() {
        notSamePasswd = "Password doesn't match";
        result = "";
      });
      developer.log(
          'checkSamePassword USER ERROR : $notSamePasswd : ${newPasswordController.text} ${repeatNewPasswordController.text}');
      return false;
    }
    setState(() {
      notSamePasswd = "";
      result = "";
    });
    return true;
  }

  Future<void> putModifyPassword() async {
    if (!checkSamePassword()) {
      return;
    }
    Response response = await ApiRequest.put(
      apiRoute, jsonEncode(<String, dynamic>{
        'oldpassword': oldPasswordController.text,
        'password': newPasswordController.text,
      }),
    );
    if (response.statusCode == 504) {
      setState(() {
        result = response.body;
      });
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "putModifyPassword FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      setState(() {
        result = '${responseData['error']}';
      });
      return;
    }
    Fluttertoast.showToast(
      msg: responseData['success'],
      toastLength: Toast.LENGTH_SHORT,
      timeInSecForIosWeb: 1,
      backgroundColor: const Color.fromARGB(255, 71, 71, 71),
      textColor: Colors.white,
      fontSize: 24.0
    );
    Navigator.pop(Globaldata.myContext);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            EditPasswordTopBar(),
            SizedBox(width: 16),
            MyDivider(),
            Text(
              result,
              style: TextStyle(
                  fontSize: 20, fontWeight: FontWeight.bold, color: Colors.red),
            ),
            LeftWidget(
              child: AutoSizeText(
                "Current password",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
                hint: 'Password',
                controller: oldPasswordController,
                obscureText: true,
                padding: EdgeInsets.symmetric(horizontal: 15)),
            Padding(padding: EdgeInsets.only(top: 30)),
            LeftWidget(
              child: AutoSizeText(
                "New password",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
                hint: 'Password',
                controller: newPasswordController,
                obscureText: true,
                padding: EdgeInsets.symmetric(horizontal: 15)),
            Padding(padding: EdgeInsets.only(top: 20)),
            LeftWidget(
              child: AutoSizeText(
                "Confirm new password",
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: Colors.black,
                ),
                maxLines: 1,
              ),
            ),
            CustomClassicForm(
                hint: 'Password',
                controller: repeatNewPasswordController,
                obscureText: true,
                padding: EdgeInsets.symmetric(horizontal: 15)),
            LeftWidget(
              child: AutoSizeText(
                notSamePasswd,
                style: TextStyle(
                  fontSize: 25,
                  fontWeight: FontWeight.bold,
                  color: const Color.fromARGB(255, 255, 0, 0),
                ),
                maxLines: 1,
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
            CustomButton(
              backgroundColor: Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText(
                "Change",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: Colors.white),
              ),
              onPressed: putModifyPassword,
              size: Size(double.infinity, 55),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
    );
  }
}

class EditPasswordTopBar extends StatelessWidget {
  const EditPasswordTopBar({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {
              Navigator.pop(Globaldata.myContext);
            },
            icon: Icon(Icons.arrow_back, size: 32),
          ),
        ),
        Center(
          child: AutoSizeText(
            "Change password",
            style: TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.w900,
              color: Colors.black,
            ),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}
