import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:convert';
import 'dart:developer' as developer;
import "../auth/webviewAuth/connectServiceAuth.dart";

import "../custom_widget.dart";
import "login_page.dart";
import "../globalData.dart";
import '../apiCall/apiRequest.dart';

class SignUp extends StatelessWidget {
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
  _APICallState createState() => _APICallState();
}

class _APICallState extends State<APICall> {
  final TextEditingController emailController = TextEditingController();
  final TextEditingController passwordController = TextEditingController();
  String result = '';
  String apiRoute = 'register';

  @override
  void dispose() {
    emailController.dispose();
    passwordController.dispose();
    super.dispose();
  }

  Future<void> postRegisterData() async {
    final response = await ApiRequest.post(apiRoute, jsonEncode(<String, dynamic>{'email': emailController.text,'password': passwordController.text,}),);
    if (response.statusCode == 504) {
      setState(() {
        result = response.body;
      });
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "postRegisterData FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      setState(() {
        result = '${responseData['error']}';
      });
      return;
    }
    Navigator.push(
      Globaldata.myContext,
      MaterialPageRoute(builder: (context) => LoginPage()),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: SingleChildScrollView(
          child: Column(
            children: [
              Padding(padding: EdgeInsets.only(top: 30)),
              Text(
                "AREA",
                style: TextStyle(fontSize: 25, fontWeight: FontWeight.bold),
              ),
              Text(
                "Sign up",
                style: TextStyle(fontSize: 40, fontWeight: FontWeight.bold),
              ),
              Text(
                result,
                style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Colors.red),
              ),
              CustomClassicForm(
                hint: 'Email',
                controller: emailController,
              ),
              CustomClassicForm(
                hint: 'Password',
                obscureText: true,
                controller: passwordController,
              ),
              SizedBox(height: 20),
              CustomButton(
                backgroundColor: Color.fromARGB(255, 34, 34, 34),
                text: AutoSizeText(
                  "Get started",
                  style: TextStyle(
                      fontSize: 30,
                      fontWeight: FontWeight.bold,
                      color: Colors.white),
                ),
                onPressed: postRegisterData,
                size: Size(double.infinity, 70),
                padding: EdgeInsets.symmetric(horizontal: 15.0),
              ),
              Padding(padding: EdgeInsets.only(top: 15)),
              MyDivider(
                text: AutoSizeText(
                  "or",
                  style: TextStyle(fontSize: 16, color: Colors.grey),
                  maxLines: 1,
                ),
                textPadding: EdgeInsets.symmetric(horizontal: 8.0),
              ),
              Padding(padding: EdgeInsets.only(top: 5)),
              AuthButton(
                text: "Continue with Spotify",
                color: Color.fromARGB(255, 30, 215, 96),
                imageUrl: 'assets/icon/Spotify.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Spotify");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              AuthButton(
                text: "Continue with Discord",
                color: Color(0xff7289da),
                imageUrl: 'assets/icon/Discord.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Discord");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              AuthButton(
                text: "Continue with Github",
                color: Color(0xff24292e),
                imageUrl: 'assets/icon/Github.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Github");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              AuthButton(
                text: "Continue with Gitlab",
                color: Color(0xfffc6d27),
                imageUrl: 'assets/icon/Gitlab.webp',
                onPressed: () {ConnectServiceAuthWebview.connectService("login", "Gitlab");},
                imageHeight: 30,
                imageWidth: 30,
              ),
              Padding(padding: EdgeInsets.only(top: 30)),
              CustomTextButton(
                text: Text(
                  "Already on AREA? Log in here",
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.bold,
                    color: Colors.black,
                  ),
                ),
                onPressed: () {
                  Navigator.push(
                    Globaldata.myContext,
                    MaterialPageRoute(builder: (context) => LoginPage()),
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
